package repository

import (
	"context"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	cmsql "gitlab.senseauto.com/apcloud/library/common-go/client/sqldb"
)

type CoreDB struct {
	MetricInfo model.CollectorMetricInfoModel
}

type CoreRepository struct {
	ctx    context.Context
	client *cmsql.Client
	DB     *CoreDB
}

func (r *CoreRepository) getMigrateTables() []interface{} {
	return []interface{}{
		&dao.CollectorMetricInfo{},
	}
}

func (r *CoreRepository) init(ctx context.Context, client *cmsql.Client) error {
	r.ctx = ctx
	r.client = client
	r.DB = &CoreDB{
		MetricInfo: model.NewCollectorMetricInfoModel(client.DB),
	}
	return nil
}
