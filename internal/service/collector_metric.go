package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	auth_filter "gitlab.senseauto.com/apcloud/library/common-go/auth/filter"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/repository"
	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
)

type ContextKeyType string

const (
	MetricTypeStartFlag  = "start_flag"
	MetricTypeEndFlag    = "end_flag"
	MetricTypeRuntimeNum = "runtime_num"

	MetricTypeProcessCost = "process_cost"
	MetricTypeSaveCost    = "save_cost"
	MetricTypeInitCost    = "init_cost"
	MetricTypeGetIDCost   = "get_id_cost"

	ContextKeyMetricInfoCollector ContextKeyType = "metric_info_collector"
	ContextKeyMetricInfoMethod    ContextKeyType = "metric_info_method"
	ContextKeyMetricInfoID        ContextKeyType = "metric_info_id"

	MetricMethodExecute = "execute"
	MetricMethodRetry   = "retry"

	MetricCollectorStatusUndone  = -1
	MetricCollectorStatusSuccess = 0
	MetricCollectorStatusFail    = 1
	MetricCollectorStatusTimeout = 2
)

// 看下现有的ticker

type CollectorMetricInfo struct {
	CollectorID string
	Method      string
	ID          uint64
	Err         error
}

func WithMetricCollector(ctx context.Context, collector Collector) context.Context {
	return context.WithValue(ctx, ContextKeyMetricInfoCollector, collector)
}
func WithMetricMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, ContextKeyMetricInfoMethod, method)
}
func WithMetricID(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, ContextKeyMetricInfoID, id)
}

func WithMetricCollectorAndMethod(ctx context.Context, collector Collector, method string) context.Context {
	return WithMetricMethod(WithMetricCollector(ctx, collector), method)
}

func contextValue[T any](ctx context.Context, key any) T {
	var t T
	v := ctx.Value(key)
	if v != nil {
		if r, ok := v.(T); ok {
			return r
		}
	}
	return t
}
func GetMetricInfo(ctx context.Context) *CollectorMetricInfo {
	info := &CollectorMetricInfo{}
	collector := contextValue[Collector](ctx, ContextKeyMetricInfoCollector)
	if collector != nil {
		info.CollectorID = collector.GetID()
	}
	method := contextValue[string](ctx, ContextKeyMetricInfoMethod)
	info.Method = method

	id := contextValue[uint64](ctx, ContextKeyMetricInfoID)
	info.ID = id
	return info
}

func CollectorMetricInit(ctx context.Context, coreDB *repository.CoreDB, defaultCollector Collector, metricMethod, metricType string, tt time.Time) (context.Context, error) {
	// 检查是否有正在执行的collector
	syncInfoService := coreDB.MetricInfo
	collectorID := defaultCollector.GetID()
	filter := model.MetricInfoFilter{
		CollectorID: &auth_filter.StringField{Eq: &collectorID}}
	record, err := syncInfoService.FindRecordByCreateTime(filter)
	if err != nil {
		return ctx, err
	}
	if record.Status == MetricCollectorStatusUndone {
		return ctx, errors.New("current collector not end!")
	}

	ctx = WithMetricCollectorAndMethod(ctx, defaultCollector, metricMethod)
	ctx = CollectorMonitorReport(ctx, coreDB, metricType, int64(time.Since(tt)), nil)
	return ctx, nil
}

func CollectorMonitorReport(ctx context.Context, coreDB *repository.CoreDB, metricType string, cost int64, err error) context.Context {
	logger := cmlog.Extract(ctx)
	info := GetMetricInfo(ctx)
	info.Err = err
	ctx, err = doCollectorMonitorReport(ctx, coreDB, info, metricType, cost)
	if err != nil {
		logger.Errorf("do collector monitor err :%s", err)
	}
	return ctx
}

func doCollectorMonitorReport(ctx context.Context, coreDB *repository.CoreDB, info *CollectorMetricInfo, metricType string, metricValue int64) (context.Context, error) {
	var err error
	cost := time.Duration(metricValue)
	syncInfoService := coreDB.MetricInfo
	if metricType == MetricTypeStartFlag || metricType == MetricTypeEndFlag {
		if info.ID == 0 {
			createAt := time.Now().Add(cost)
			insertRecord := &dao.CollectorMetricInfo{
				CollectorID: info.CollectorID,
				Method:      info.Method,
				CreateAt:    createAt.UnixMilli(),
				Status:      MetricCollectorStatusUndone,
			}
			insertID, err := syncInfoService.InsertRecord(insertRecord)
			context := WithMetricID(ctx, insertID)
			return context, err
		} else {
			status := MetricCollectorStatusSuccess
			extraInfo := ""
			if info.Err != nil {
				status = MetricCollectorStatusFail
				extraInfo = info.Err.Error()
			}
			updateRecord := &dao.CollectorMetricInfo{
				ID:        info.ID,
				Duration:  cost.Seconds(),
				Status:    int32(status),
				ExtraInfo: extraInfo,
			}
			err = syncInfoService.UpdateRecordByID(updateRecord)
		}
	}
	if metricType == MetricTypeRuntimeNum {
		collectorID := info.CollectorID
		filter := model.MetricInfoFilter{
			CollectorID: &auth_filter.StringField{Eq: &collectorID}}
		record, err := syncInfoService.FindRecordByCreateTime(filter)
		if err != nil {
			return ctx, err
		}
		if record.ID != 0 {
			record.RuntimeNum = int32(metricValue)
			if info.Err != nil {
				extraInfo := info.Err.Error()
				record.ExtraInfo += extraInfo
			}
			_ = syncInfoService.UpdateRecordByLatestID(record)
		}
	}
	return ctx, err
}
