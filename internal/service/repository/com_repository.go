package repository

import (
	"context"

	"github.com/fabxu/datacollector-service/internal/model"
	"github.com/fabxu/datacollector-service/internal/model/dao"
	cmsql "github.com/fabxu/lib/client/sqldb"
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
