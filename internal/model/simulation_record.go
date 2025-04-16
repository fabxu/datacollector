package model

import (
	"context"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	cmclient "gitlab.senseauto.com/apcloud/library/common-go/client"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SimulationRecordModel interface {
	BatchInsertSimulationRecord(ctx context.Context, records []*dao.SimulationRecord) error
	GetSimulationRecordByStatus(ctx context.Context, status int32) ([]*dao.SimulationRecord, error)
	BatchUpdateSimulationRecord(ctx context.Context, records []*dao.SimulationRecord) error
	UpdateSimulationRecord(ctx context.Context, record *dao.SimulationRecord) error
	BatchUpsertSimulationRecord(ctx context.Context, records []*dao.SimulationRecord) error
	UpsertSimulationRecord(ctx context.Context, record *dao.SimulationRecord) error
}

type simulationRecordModelImpl struct {
	db *gorm.DB
}

func NewSimulationRecordModel() SimulationRecordModel {
	return &simulationRecordModelImpl{db: cmclient.SQLDB.DB}
}

func (m *simulationRecordModelImpl) BatchUpsertSimulationRecord(ctx context.Context, records []*dao.SimulationRecord) error {
	if len(records) == 0 {
		return nil
	}
	for _, record := range records {
		_ = m.UpsertSimulationRecord(ctx, record)
	}
	return nil
}
func (m *simulationRecordModelImpl) UpsertSimulationRecord(ctx context.Context, record *dao.SimulationRecord) error {
	tx := cmclient.SQLDB.Extract(ctx, m.db)
	var count int64
	db := tx.Model(&dao.SimulationRecord{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("job_id = ?", record.JobID)

	if err := db.Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return m.UpdateSimulationRecord(ctx, record)
	}

	return tx.Create(record).Error
}

func (m *simulationRecordModelImpl) BatchInsertSimulationRecord(ctx context.Context, records []*dao.SimulationRecord) error {
	tx := cmclient.SQLDB.Extract(ctx, m.db)
	return tx.CreateInBatches(records, len(records)).Error
}

func (m *simulationRecordModelImpl) GetSimulationRecordByStatus(ctx context.Context, status int32) ([]*dao.SimulationRecord, error) {
	tx := cmclient.SQLDB.Extract(ctx, m.db)

	var records []*dao.SimulationRecord

	err := tx.Model(&dao.SimulationRecord{}).
		Where("status = ?", status).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (m *simulationRecordModelImpl) BatchUpdateSimulationRecord(ctx context.Context, records []*dao.SimulationRecord) error {

	if len(records) == 0 {
		return nil
	}
	for _, record := range records {
		_ = m.UpdateSimulationRecord(ctx, record)
	}
	return nil
}

func (m *simulationRecordModelImpl) UpdateSimulationRecord(ctx context.Context, record *dao.SimulationRecord) error {

	return nil
}
