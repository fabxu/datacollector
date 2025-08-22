package simulation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/fabxu/datacollector-service/internal/lib/constant"
	"github.com/fabxu/datacollector-service/internal/model"
	"github.com/fabxu/datacollector-service/internal/model/dao"
	"github.com/fabxu/datacollector-service/internal/service"
	"github.com/fabxu/datacollector-service/internal/service/monitor"
	"github.com/fabxu/datacollector-service/internal/service/util"
	cmconfig "github.com/fabxu/lib/config"
	cmlog "github.com/fabxu/log"
)

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

	return nil, nil, nil
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
