package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/lib/constant"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/repository"
	cmhttp "gitlab.senseauto.com/apcloud/library/common-go/client/http"
	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
)

const (
	statusUnInit = iota
	statusRunning
	statusClosing

	defaultCollectorTickTime   = 1
	defaultMinCollectorTTT     = 5
	defaultCollectorReportTime = 60
)

type Options struct {
	tickerTime int64  // collector 基础时钟
	cronRule   string // cronjob规则

	ttt int64 // 自身tick的周期，单位秒
	ttr int64 // 上报周期，单位秒
}

func (o *Options) CheckAndReset() error {
	o.tickerTime = defaultCollectorTickTime
	o.ttr = defaultCollectorReportTime

	if o.cronRule != "" {
		if _, err := cron.ParseStandard(o.cronRule); err != nil {
			return err
		}
	}

	if o.ttt != 0 {
		if o.ttt < defaultMinCollectorTTT {
			o.ttt = defaultMinCollectorTTT
		}
	}

	return nil
}

type Option func(*Options)

type TickFunc func(CollectorID string) error

func WithTickerTime(tickerTime int64) Option {
	return func(o *Options) {
		o.tickerTime = tickerTime
	}
}
func WithTTT(ttt int64) Option {
	return func(o *Options) {
		o.ttt = ttt
	}
}
func WithCronRule(cronRule string) Option {
	return func(o *Options) {
		o.cronRule = cronRule
	}
}

type collectorRuntime struct {
	ctx            context.Context
	collector      Collector
	coreRepository *repository.CoreRepository
	httpClient     *cmhttp.Wrapper
	opts           *Options
	doOnce         sync.Once

	runtimeNum int64
	tickStatus int32
	quitCh     chan struct{}
	ticker     *time.Ticker
	tickerFunc TickFunc
	cronStatus int32
	cron       *cron.Cron

	lastTickTime   int64
	lastReportTime int64
	lastCronTime   int64
}

func (cr *collectorRuntime) serve(ctx context.Context) {
	logger := cmlog.Extract(ctx)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered from %v", err)
				cr.close()
			}
		}()
		for cr.tickStatus == statusRunning {
			select {
			case <-cr.quitCh:
				logger.Debugf("collector【%s】 runtime stopped,status :%s", cr.collector.GetID(), cr.tickStatus)
				cr.runtimeNum--
				return
			case t := <-cr.ticker.C:
				cr.tick(t)
			}
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered from %v", err)
				cr.cron.Stop()
				cr.cronStatus = statusClosing
			}
		}()
		if cr.cronStatus == statusRunning {
			cr.runtimeNum++
			_, _ = cr.cron.AddFunc(cr.opts.cronRule, func() {
				cr.lastCronTime = time.Now().Unix()
				cr.executeCron(ctx)
			})
			cr.cron.Start()
		}
	}()
}

func (cr *collectorRuntime) executeCron(ctx context.Context) {
	logger := cmlog.Extract(ctx)
	cli := cr.httpClient.GetClient(constant.AppName)
	request := cli.R()
	data := make(map[string]string)
	data["id"] = cr.collector.GetID()

	resp, err := request.SetHeader("Content-Type", "application/json;charset=UTF-8").
		SetBody(&data).
		Post(constant.AppCollectorPath)
	if err == nil {
		var rspBody dao.Response[interface{}]

		if err = json.Unmarshal(resp.Body(), &rspBody); err != nil {
			logger.Error(err.Error())
			logger.Error(string(resp.Body()))
		} else if rspBody.Code != constant.HTTPCodeSuccessCollector {
			logger.Errorf("collector(%s) execute code(%s),resp(%s)", cr.collector.GetID(), rspBody.Code, rspBody.Message)
		}
	} else {
		logger.Errorf("collector(%s) execute err(%s),resp(%s)", cr.collector.GetID(), err.Error(), string(resp.Body()))
	}

}

func (cr *collectorRuntime) tick(tt time.Time) {

	var err error = nil
	logger := cmlog.Extract(cr.ctx)
	t := tt.Unix()
	if cr.opts.ttt > 0 && t > cr.opts.ttt+cr.lastTickTime {
		err = cr.tickerFunc(cr.collector.GetID())
		if err != nil {
			logger.Errorf("collector(%s) tick err(%s),elapsed(%s)", cr.collector.GetID(), time.Since(tt))
		}
		cr.lastTickTime = t
		cr.doOnce.Do(func() {
			cr.runtimeNum++
		})
	}

	if cr.lastCronTime > 0 {
		lastCronTimeUnix := time.Unix(cr.lastCronTime, 0)
		schedule, _ := cron.ParseStandard(cr.opts.cronRule)
		nextTime := schedule.Next(lastCronTimeUnix)
		interval := nextTime.Sub(lastCronTimeUnix)

		if tt.Sub(lastCronTimeUnix) > 2*interval {
			cr.doOnce.Do(func() {
				cr.runtimeNum--
			})
		}
	}

	if cr.opts.ttr > 0 && t > cr.opts.ttr+cr.lastReportTime {
		cr.ctx = WithMetricCollector(cr.ctx, cr.collector)
		CollectorMonitorReport(cr.ctx, cr.coreRepository.DB, MetricTypeRuntimeNum, cr.runtimeNum, err)
		cr.lastReportTime = t
	}

}

func (cr *collectorRuntime) close() {
	cr.runtimeNum--
	cr.tickStatus = statusClosing
	cr.quitCh <- struct{}{}
	close(cr.quitCh)
	cr.ticker.Stop()
}
