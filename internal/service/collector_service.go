package service

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/robfig/cron/v3"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/lib/constant"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/util"
	cmhttp "gitlab.senseauto.com/apcloud/library/common-go/client/http"

	"github.com/golang/protobuf/ptypes/empty"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/monitor"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/repository"
	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
	dc_api "gitlab.senseauto.com/apcloud/library/proto/api/datacollector-service/v1"
)

type CollectorService struct {
	dc_api.UnsafeCollectServiceServer
	Collectors     map[string]Collector
	Monitor        *monitor.Monitor
	ComRepository  *repository.ComRepository
	CoreRepository *repository.CoreRepository

	Ctx context.Context
}

func (s *CollectorService) RegisterCollector(collector Collector) error {
	if collector == nil {
		return errors.New("collector is nil")
	}

	if _, ok := s.Collectors[collector.GetID()]; ok {
		return errors.New("Collector is dulplicate! Collector : " + collector.GetID())
	} else {
		s.Collectors[collector.GetID()] = collector
	}

	return nil
}

func (s *CollectorService) GetMigrateTables() []interface{} {
	tableMap := make(map[string]interface{})
	result := make([]interface{}, 0)
	for _, value := range s.Collectors {
		tables := value.GetMigrateTables()
		for _, item := range tables {
			key := reflect.TypeOf(item).String()
			if _, exist := tableMap[key]; !exist {
				tableMap[key] = item
				result = append(result, item)
			}
		}
	}
	return result
}

func (s *CollectorService) Init() error {
	var msg string
	var err error
	s.ComRepository, s.CoreRepository, err = repository.CreateRepository(s.Ctx)
	if err == nil {
		s.Monitor.Init(s.Ctx)
		collectorCtx := createContext(s)

		for _, collector := range s.Collectors {
			err = collector.Init(collectorCtx)
			if err != nil {
				msg = msg + err.Error() + ","
			}
		}
	} else {
		msg = err.Error()
	}

	if len(msg) > 0 {
		return errors.New(msg)
	}

	return nil
}

func (s *CollectorService) Tick() error {

	var msg string
	var err error
	for _, collector := range s.Collectors {

		opts := &Options{}
		tickFunc, opt := collector.GetTickOpts()
		for _, o := range opt {
			o(opts)
		}
		if err = opts.CheckAndReset(); err != nil {
			msg = msg + collector.GetID() + err.Error() + ","
			continue
		}
		httpCfg := cmhttp.Config{
			Host:               util.GetProjectCollectorURL(),
			Timeout:            30,
			InsecureSkipVerify: true,
		}
		t := time.Now().Unix()
		cr := &collectorRuntime{
			ctx:            s.Ctx,
			collector:      collector,
			coreRepository: s.CoreRepository,
			opts:           opts,
			httpClient:     cmhttp.New(s.Ctx, constant.AppName, httpCfg),
			tickStatus:     statusRunning,
			ticker:         time.NewTicker(time.Duration(opts.tickerTime) * time.Second),
			cronStatus:     statusUnInit,
			quitCh:         make(chan struct{}),
			runtimeNum:     0,
			lastTickTime:   t,
			lastReportTime: t,
		}

		if opts.ttt != 0 {
			cr.tickerFunc = tickFunc
		}
		if opts.cronRule != "" {
			cr.cronStatus = statusRunning
			cr.cron = cron.New()
		}
		cr.serve(s.Ctx)
	}

	if len(msg) > 0 {
		return errors.New(msg)
	}
	return nil
}

func (s *CollectorService) ExecuteCollect(
	ctx context.Context,
	req *dc_api.CollectRequest,
) (*dc_api.CollectResponse, error) {
	var err error
	var value interface{}
	var msg *monitor.MonitorMsg
	logger := cmlog.Extract(ctx)
	startTime := time.Now()
	ctx, err = CollectorMetricInit(ctx, s.CoreRepository.DB, s.Collectors[req.Info.Id], MetricMethodExecute, MetricTypeStartFlag, startTime)
	if err != nil {
		return &dc_api.CollectResponse{
			Code:    "",
			Message: err.Error(),
		}, nil
	}
	defer func() {
		if err != nil {
			logger.Errorf("【execute】 err:%s", err)
		}
		CollectorMonitorReport(ctx, s.CoreRepository.DB, MetricTypeEndFlag, int64(time.Since(startTime)), err)
	}()

	if collector, ok := s.Collectors[req.Info.Id]; ok {
		value, msg, err = collector.Process(ctx, req.Info.Type, req.Info.Arg)
		if err == nil {
			if err = collector.Save(ctx, value); err != nil {
				logger.Error("Collect save failed! err : " + err.Error())
			}
			logger.Info("finish collect process!")
		} else {
			logger.Error(err)
		}

		if msg != nil {
			msg.RetryMsg.ID = req.Info.Id
			s.Monitor.HandleMonitorMsg(msg)
		}
	} else {
		logger.Error("No such collector, id : " + req.Info.Id)
	}

	return &dc_api.CollectResponse{}, nil
}

func (s *CollectorService) Retry(
	ctx context.Context, empty *empty.Empty) (*dc_api.CollectResponse, error) {
	jobs, err := s.Monitor.GetRetryJobs()
	logger := cmlog.Extract(ctx)

	if err == nil {
		count := len(jobs)

		if count > 0 {
			for _, job := range jobs {
				if collector, ok := s.Collectors[job.ID]; ok {
					go func() {
						value, msg, err := collector.Process(ctx, job.MsgType, job.Arg)

						if err == nil {
							if err := collector.Save(ctx, value); err != nil {
								logger.Error("Collector save failed! error : " + err.Error())
							}
						}

						if msg != nil {
							msg.RetryMsg.ID = job.ID
							s.Monitor.HandleMonitorMsg(msg)
						}
					}()
				} else {
					logger.Error(err)
				}
			}
		}
	}

	return &dc_api.CollectResponse{}, nil
}
