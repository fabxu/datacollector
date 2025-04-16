package model

import (
	"errors"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	auth_filter "gitlab.senseauto.com/apcloud/library/common-go/auth/filter"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CollectorSyncInfoModel interface {
	InsertRecord(record *dao.CollectorSyncInfo) (uint64, error)
	FindRecordByID(filter SyncInfoFilter) (*dao.CollectorSyncInfo, error)
}

type SyncInfoFilter struct {
	ID          *auth_filter.UInt64Field
	CollectorID *auth_filter.StringField
}

type collectorSyncInfoModelImpl struct {
	db *gorm.DB
}

func NewCollectorSyncInfoModel(db *gorm.DB) CollectorSyncInfoModel {
	return &collectorSyncInfoModelImpl{db: db}
}

// FindRecordByID implements CollectorSyncInfoModel.
func (m *collectorSyncInfoModelImpl) FindRecordByID(filter SyncInfoFilter) (*dao.CollectorSyncInfo, error) {
	tx := m.db.Model(&dao.CollectorSyncInfo{})

	tx = filter.ID.ToAndWhere(tx, "id")
	tx = filter.CollectorID.ToAndWhere(tx, "collector_id")
	tx.Order("last_sync_time DESC")
	syncInfo := dao.CollectorSyncInfo{}

	result := tx.First(&syncInfo)
	if result.Error != nil {
		// 如果记录没有找到，返回空记录，不报错
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &syncInfo, nil
		}

		return nil, result.Error
	}

	return &syncInfo, nil
}

func (m *collectorSyncInfoModelImpl) InsertRecord(record *dao.CollectorSyncInfo) (uint64, error) {
	tx := m.db
	tx = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(record)
	return record.ID, tx.Error
}
