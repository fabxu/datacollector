package model

import (
	"github.com/fabxu/datacollector-service/internal/model/dao"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CollectorSyncInfoModel interface {
	InsertRecord(record *dao.CollectorSyncInfo) (uint64, error)
}

type collectorSyncInfoModelImpl struct {
	db *gorm.DB
}

func NewCollectorSyncInfoModel(db *gorm.DB) CollectorSyncInfoModel {
	return &collectorSyncInfoModelImpl{db: db}
}

func (m *collectorSyncInfoModelImpl) InsertRecord(record *dao.CollectorSyncInfo) (uint64, error) {
	tx := m.db
	tx = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(record)
	return record.ID, tx.Error
}
