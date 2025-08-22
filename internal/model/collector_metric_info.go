package model

import (
	"github.com/fabxu/datacollector-service/internal/model/dao"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CollectorMetricInfoModel interface {
	InsertRecord(record *dao.CollectorMetricInfo) (uint64, error)
	UpdateRecordByID(record *dao.CollectorMetricInfo) error
	UpdateRecordByLatestID(record *dao.CollectorMetricInfo) error
}

type collectorMetricInfoModelImpl struct {
	db *gorm.DB
}

func NewCollectorMetricInfoModel(db *gorm.DB) CollectorMetricInfoModel {
	return &collectorMetricInfoModelImpl{db: db}
}

func (m *collectorMetricInfoModelImpl) InsertRecord(record *dao.CollectorMetricInfo) (uint64, error) {
	tx := m.db
	tx = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(record)
	return record.ID, tx.Error
}

func (m *collectorMetricInfoModelImpl) UpdateRecordByID(record *dao.CollectorMetricInfo) error {
	tx := m.db
	if err := tx.Model(&dao.CollectorMetricInfo{}).
		Where("id = ?", record.ID).
		Updates(map[string]interface{}{
			"status":     record.Status,
			"duration":   record.Duration,
			"extra_info": record.ExtraInfo,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (m *collectorMetricInfoModelImpl) UpdateRecordByLatestID(record *dao.CollectorMetricInfo) error {

	tx := m.db
	tx = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(record)
	return tx.Error
}
