package main

import (
	"context"

	"github.com/fabxu/datacollector-service/internal/model"
	"github.com/fabxu/datacollector-service/internal/service"
	"github.com/fabxu/datacollector-service/internal/service/monitor"
	"github.com/fabxu/datacollector-service/internal/service/repository"
	"github.com/fabxu/datacollector-service/internal/service/simulation"
	"github.com/fabxu/lib/client/sqldb"
	cmlog "github.com/fabxu/log"
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
