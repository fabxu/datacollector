package main

import (
	"context"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/monitor"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/repository"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/service/simulation"
	"gitlab.senseauto.com/apcloud/library/common-go/client/sqldb"
	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
	dc_api "gitlab.senseauto.com/apcloud/library/proto/api/datacollector-service/v1"
)

var serviceInstance = service.CollectorService{}

func GetServiceInstance() dc_api.CollectServiceServer {
	return &serviceInstance
}

func RegisterService(ctx context.Context) {
	serviceInstance.Collectors = make(map[string]service.Collector)
	serviceInstance.Monitor = monitor.NewMonitor(ctx)
	serviceInstance.Ctx = ctx
	_ = serviceInstance.RegisterCollector(simulation.NewSimulationCollector(ctx))
}

func StartService(ctx context.Context) {
	logger := cmlog.Extract(ctx)

	err := serviceInstance.Init()
	if err != nil {
		logger.Errorf("collector service 【init】  failed,err:%s", err)
	}

	// err = serviceInstance.Tick()
	//  if err != nil {
	//	logger.Errorf("collector service 【tick】 failed,err:%s", err)
	//  }

}

func AutoMigrate(ctx context.Context, client *sqldb.Client) {
	model.AutoMigrate(client, serviceInstance.GetMigrateTables())
	repository.AutoMigrate(ctx)
}
