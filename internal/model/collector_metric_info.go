package model

import (
	"errors"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	auth_filter "gitlab.senseauto.com/apcloud/library/common-go/auth/filter"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CollectorMetricInfoModel interface {
	InsertRecord(record *dao.CollectorMetricInfo) (uint64, error)
	FindRecordByID(filter MetricInfoFilter) (*dao.CollectorMetricInfo, error)
	UpdateRecordByID(record *dao.CollectorMetricInfo) error
	FindRecordByCreateTime(filter MetricInfoFilter) (*dao.CollectorMetricInfo, error)
	UpdateRecordByLatestID(record *dao.CollectorMetricInfo) error
}

type MetricInfoFilter struct {
	ID          *auth_filter.UInt64Field
	CollectorID *auth_filter.StringField
}

type collectorMetricInfoModelImpl struct {
	db *gorm.DB
}

func NewCollectorMetricInfoModel(db *gorm.DB) CollectorMetricInfoModel {
	return &collectorMetricInfoModelImpl{db: db}
}

// FindRecordByID implements CollectorSyncInfoModel.
func (m *collectorMetricInfoModelImpl) FindRecordByID(filter MetricInfoFilter) (*dao.CollectorMetricInfo, error) {
	tx := m.db.Model(&dao.CollectorMetricInfo{})

	tx = filter.ID.ToAndWhere(tx, "id")
	tx = filter.CollectorID.ToAndWhere(tx, "collector_id")

	metricInfo := dao.CollectorMetricInfo{}

	result := tx.First(&metricInfo)
	if result.Error != nil {
		// 如果记录没有找到，返回空记录，不报错
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &metricInfo, nil
		}
		return nil, result.Error
	}

	return &metricInfo, nil
}
func (m *collectorMetricInfoModelImpl) FindRecordByCreateTime(filter MetricInfoFilter) (*dao.CollectorMetricInfo, error) {
	tx := m.db.Model(&dao.CollectorMetricInfo{})

	tx = filter.CollectorID.ToAndWhere(tx, "collector_id")
	tx.Order("create_at DESC")
	metricInfo := dao.CollectorMetricInfo{}
	result := tx.First(&metricInfo)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return &metricInfo, nil
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
