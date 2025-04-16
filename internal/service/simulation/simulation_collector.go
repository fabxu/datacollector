package simulation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	auth_filter "gitlab.senseauto.com/apcloud/library/common-go/auth/filter"

	cmlib "gitlab.senseauto.com/apcloud/library/common-go/lib"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/lib/constant"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/monitor"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/util"
	cmconfig "gitlab.senseauto.com/apcloud/library/common-go/config"
	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
)

var defaultStartTime = time.Date(2023, 1, 1, 0, 0, 0, 0, cmlib.GetCSTLocation())
var JobStatusMap = map[constant.WorkflowJobStatus]int{
	constant.JobFinished:   0,
	constant.JobUnfinished: 1,
	constant.JobPause:      2,
	constant.JobRunning:    3,
	constant.JobPending:    4,
}

const (
	AlertModuleGetSimulation = "GetSimulation"
	CronRule                 = "0 6 * * *"
)

type dataHolder struct {
	simData  []*dao.SimulationRecord
	syncInfo *dao.CollectorSyncInfo
}

func getSimRecordPath() (string, string) {
	return "POST", "/api/v1/sim/jobForInternal/biDashboard"
}

func (s *SimulationCollector) getSimRecord(ctx context.Context, startTime int64, endTime int64) ([]*dao.SimulationRecord, error) {
	var err error
	var simRecords []*dao.SimulationRecord
	logger := cmlog.Extract(ctx)
	method, path := getSimRecordPath()
	url := fmt.Sprintf("%s%s", s.config.SimURL, path)
	data := make(map[string]interface{})
	// data["StartTime"], data["EndTime"] = util.GetYesterDateStartAndEndTime(time.Now())
	data["StartTime"] = time.UnixMilli(startTime).Format(constant.FullTimeTemplate)
	data["EndTime"] = time.UnixMilli(endTime).Format(constant.FullTimeTemplate)

	data["Page"] = 1
	data["Pagesize"] = -1
	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	result, _ := util.DoRequest(method, url, byteData, headers)
	var rspBody dao.SimRecordResponse
	body := []byte(result)
	if err := json.Unmarshal(body, &rspBody); err != nil {
		return nil, err
	}
	if rspBody.Message != "OK" {
		return nil, errors.New(rspBody.Message)
	}
	for _, record := range rspBody.Data.Jobs {
		status := int32(JobStatusMap[record.Status])
		if _, ok := JobStatusMap[record.Status]; !ok {
			status = -1
			logger.Errorf("jobid:【%s】,status:【%d】 abnormal!", record.JobID, status)
		}
		tmpRecord := &dao.SimulationRecord{
			JobID:             record.JobID,
			Name:              record.Name,
			SpaceID:           record.SpaceID,
			SpaceType:         record.SpaceIdcn,
			SpaceTypeCategory: record.DataOrigin,
			TotalDistance:     record.TotalDistance,
			Status:            status,
			CreateTime:        util.TimeParse(record.CreateTime),
			StartTime:         util.TimeParse(record.StartTime),
			EndTime:           util.TimeParse(record.EndTime),
			Creator:           record.Creator,
		}
		simRecords = append(simRecords, tmpRecord)
	}
	return simRecords, nil
}

type SimulationCollector struct {
	service.Collector
	*service.CollectorContext
	simulationRecordModel model.SimulationRecordModel
	config                simConfig
}

type simConfig struct {
	SimURL string `mapstructure:"sim_url"`
}

func (s *SimulationCollector) GetID() string {
	return constant.IDCollectorSimulation
}

func (s *SimulationCollector) GetMigrateTables() []interface{} {
	return []interface{}{
		&dao.SimulationRecord{},
	}
}

func (s *SimulationCollector) GetTickOpts() (tickFunc service.TickFunc, options []service.Option) {

	options = []service.Option{
		service.WithCronRule(CronRule),
		// service.WithTTT(3),
	}

	tickFunc = func(CollectorID string) error {
		// fmt.Println("this is a tick function", CollectorID)
		return nil
	}

	return
}

// init implements Collector.
func (s *SimulationCollector) Init(ctx *service.CollectorContext) error {
	s.CollectorContext = ctx
	logger := cmlog.Extract(s.Ctx)
	if err := cmconfig.Global().UnmarshalKey(constant.CfgSimulation, &s.config); err != nil {
		logger.Panic(err)
	}
	return nil
}

// process implements Collector.
func (s *SimulationCollector) Process(ctx context.Context, msgType int32, arg string) (interface{}, *monitor.MonitorMsg, error) {
	var msg *monitor.MonitorMsg = nil
	var err error
	collectorID := s.GetID()
	currentTime := time.Now()
	current := cmlib.ToMilliSeconds(currentTime)
	filter := model.SyncInfoFilter{
		CollectorID: &auth_filter.StringField{Eq: &collectorID}}
	holder := dataHolder{}
	if holder.syncInfo, err = s.Repository.DB.SyncInfo.FindRecordByID(filter); err == nil {
		if holder.syncInfo.CreateAt == 0 {
			holder.syncInfo.CollectorID = s.GetID()
			holder.syncInfo.CreateAt = current
		}
		lastUpdateTime := holder.syncInfo.LastSyncTime
		if lastUpdateTime == 0 {
			lastUpdateTime = cmlib.ToMilliSeconds(defaultStartTime)
		}
		holder.simData, err = s.getSimRecord(ctx, lastUpdateTime, current)
	}

	if err != nil {
		alertMsg := make([]*monitor.AlertMsg, 1)
		alertMsg[0] = &monitor.AlertMsg{
			AlertType: monitor.AlertService,
			Service:   s.GetID(),
			Module:    AlertModuleGetSimulation,
			Time:      currentTime.Format(constant.FullTimeTemplate),
			Msg:       err.Error(),
		}
		msg = &monitor.MonitorMsg{
			Retry:    false,
			Alert:    true,
			AlertMsg: alertMsg,
		}
	} else {
		holder.syncInfo.LastSyncTime = current
	}

	return &holder, msg, err
}

// save implements Collector.
func (s *SimulationCollector) Save(ctx context.Context, param interface{}) error {
	data := param.(*dataHolder)
	if len(data.simData) > 0 {
		_ = s.simulationRecordModel.BatchUpsertSimulationRecord(ctx, data.simData)
	}
	_, _ = s.Repository.DB.SyncInfo.InsertRecord(data.syncInfo)
	return nil
}

func NewSimulationCollector(ctx context.Context) service.Collector {
	collector := &SimulationCollector{
		simulationRecordModel: model.NewSimulationRecordModel(),
	}
	return collector
}
