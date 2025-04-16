package repository

import (
	"context"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	cmsql "gitlab.senseauto.com/apcloud/library/common-go/client/sqldb"
)

type ComDB struct {
	SyncInfo model.CollectorSyncInfoModel
	FileInfo model.FileInfoModel
}

type ComRepository struct {
	ctx    context.Context
	client *cmsql.Client
	DB     *ComDB
}

func (r *ComRepository) getMigrateTables() []interface{} {
	return []interface{}{
		&dao.CollectorSyncInfo{},
		&dao.FileInfo{},
	}
}

func (r *ComRepository) init(ctx context.Context, client *cmsql.Client) error {
	r.ctx = ctx
	r.client = client
	r.DB = &ComDB{
		SyncInfo: model.NewCollectorSyncInfoModel(client.DB),
		FileInfo: model.NewFileInfoModel(client.DB),
	}
	return nil
}
