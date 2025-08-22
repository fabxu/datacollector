package service

import (
	"context"

	"github.com/fabxu/datacollector-service/internal/service/monitor"
	"github.com/fabxu/datacollector-service/internal/service/repository"
)

type Collector interface {
	GetMigrateTables() []interface{}
	Init(ctx *CollectorContext) error
	GetID() string
	Process(ctx context.Context, msgType int32, arg string) (interface{}, *monitor.MonitorMsg, error)
	Save(ctx context.Context, data interface{}) error
	GetTickOpts() (tick TickFunc, options []Option)
}

type CollectorContext struct {
	Ctx        context.Context
	Repository *repository.ComRepository
}

func createContext(service *CollectorService) *CollectorContext {
	return &CollectorContext{
		Ctx:        service.Ctx,
		Repository: service.ComRepository,
	}
}
