package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/fabxu/datacollector-service/internal/lib/constant"
	cmclient "github.com/fabxu/lib/client"
	cmredis "github.com/fabxu/lib/client/redis"
	cmconfig "github.com/fabxu/lib/config"
	cmlog "github.com/fabxu/log"
	"github.com/go-redis/redis/v8"
)

type retryMsgType int

const (
	typeRetryMsgSubmit retryMsgType = iota
	typeRetryMsgFetch
	typeRetryMsgFetchReply
)

const (
	monitorTickInterval = 500
	monitorCallTimeout  = 5
)

type RetryMsg struct {
	ID       string
	Interval uint32
	MsgType  int32
	Arg      string
}

type handleMsg struct {
	msgType retryMsgType
	data    interface{}
	err     error
}

func (m RetryMsg) equal(msg *RetryMsg) bool {
	if m.ID == msg.ID && m.MsgType == msg.MsgType && m.Arg == msg.Arg {
		return true
	}

	return false
}

type MonitorMsg struct {
	Retry    bool
	RetryMsg RetryMsg
	Alert    bool
	AlertMsg []*AlertMsg
}

type MonitorConfig struct {
	Interval uint32
	//nolint:stylecheck
	Cache_size uint32
}

type Monitor struct {
	ctx      context.Context
	config   MonitorConfig
	receiver chan *handleMsg
	sender   chan *handleMsg
	alert    *Alert
}

func NewMonitor(ctx context.Context) *Monitor {
	monitor := Monitor{ctx: ctx, config: MonitorConfig{}, alert: newAlert(ctx)}
	logger := cmlog.Extract(ctx)

	if err := cmconfig.Global().UnmarshalKey(constant.CfgMonitor, &monitor.config); err != nil {
		logger.Panic(err)
	}

	monitor.receiver = make(chan *handleMsg, monitor.config.Cache_size)
	monitor.sender = make(chan *handleMsg, monitor.config.Cache_size)

	return &monitor
}

func (m *Monitor) task() {
	ticker := time.NewTicker(monitorTickInterval * time.Millisecond)

	for {
		select {
		case msg := <-m.receiver:
			switch msg.msgType {
			case typeRetryMsgSubmit:
				data := msg.data.(RetryMsg)
				m._saveRetryJob(data)
			case typeRetryMsgFetch:
				retryMsg, err := m.getRetryJobs()
				m.sender <- &handleMsg{msgType: typeRetryMsgFetchReply, data: retryMsg, err: err}
			default:
			}
		case <-ticker.C:
		}
	}
}

func (m *Monitor) Init(ctx context.Context) {
	redisCfg := cmredis.Config{}
	logger := cmlog.Extract(ctx)

	if err := cmconfig.Global().UnmarshalKey(constant.CfgRedis, &redisCfg); err != nil {
		logger.Panic(err)
	}

	redisCfg.Single.DB = constant.CfgRedisDB
	redisCfg.Type = constant.CfgRedisType

	cmclient.Redis.Global(ctx, redisCfg)
	m.alert.init(ctx)

	go m.task()
}

func (m *Monitor) HandleMonitorMsg(msg *MonitorMsg) {
	if msg.Alert {
		m.alert.sendAlert(msg.AlertMsg)
	}

	if msg.Retry {
		m.receiver <- &handleMsg{msgType: typeRetryMsgSubmit, data: msg.RetryMsg}
	}
}

func (m *Monitor) GetRetryJobs() ([]RetryMsg, error) {
	m.receiver <- &handleMsg{msgType: typeRetryMsgFetch}
	select {
	case msg := <-m.sender:
		if msg.data != nil {
			data := msg.data.([]RetryMsg)
			return data, msg.err
		}

		return nil, msg.err
	case <-time.After(monitorCallTimeout * time.Second):
		return nil, errors.New("get retry job timeout")
	}
}

func (m *Monitor) getRetryJobs() ([]RetryMsg, error) {
	logger := cmlog.Extract(m.ctx)
	value, err := cmclient.Redis.Get(m.ctx, constant.KeyRedisRetryJobs).Bytes()
	result := make([]RetryMsg, 0)

	if err == nil {
		var jobs []RetryMsg
		err = json.Unmarshal(value, &jobs)

		if err == nil {
			index := -1

			for i, value := range jobs {
				value.Interval--
				// find the last index where job interval is 0
				if value.Interval == 0 {
					index = i
				}
			}

			if index >= 0 {
				if index == len(jobs) {
					result = jobs

					cmclient.Redis.Del(m.ctx, constant.KeyRedisRetryJobs)
				} else {
					result = jobs[:index+1]
					jobs = jobs[index+1:]
					bytes, _ := json.Marshal(jobs)
					cmclient.Redis.Set(m.ctx, constant.KeyRedisRetryJobs, bytes, constant.DurationExpire) // 这里有问题 24s ？
				}
			}
		}
	} else if err != redis.Nil {
		logger.Panic(err)
	}

	return result, err
}

func (m *Monitor) _interal2Tick(interval uint32) uint32 {
	interval /= m.config.Interval
	if interval == 0 {
		interval++
	}

	return interval
}

func (m *Monitor) _saveRetryJob(msg RetryMsg) {
	var jobs []RetryMsg

	logger := cmlog.Extract(m.ctx)
	msg.Interval = m._interal2Tick(msg.Interval)

	value, err := cmclient.Redis.Get(m.ctx, constant.KeyRedisRetryJobs).Bytes()

	switch err {
	case redis.Nil:
		jobs = make([]RetryMsg, 1)
		jobs[0] = msg
	case nil:
		err = json.Unmarshal(value, &jobs)

		if err == nil {
			if len(jobs) == 0 || jobs[len(jobs)-1].Interval <= msg.Interval {
				jobs = append(jobs, msg)
			} else {
				var value RetryMsg

				i := 0
				flag := true

				for i, value = range jobs {
					if value.Interval >= msg.Interval {
						break
					}

					if msg.equal(&value) {
						flag = false
						break
					}

					i++
				}

				if flag {
					jobs = append(jobs, msg)
					copy(jobs[i+1:], jobs[i:])
					jobs[i] = msg
				}
			}
		} else {
			logger.Error(err)
		}
	default:
		logger.Panic(err)
	}

	if err == nil || err == redis.Nil {
		value, err = json.Marshal(jobs)
		if err == nil {
			cmclient.Redis.Set(m.ctx, constant.KeyRedisRetryJobs, value, constant.DurationExpire)
		} else {
			logger.Error(err)
		}
	}
}
