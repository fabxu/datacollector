package model

import (
	"github.com/fabxu/datacollector-service/internal/model/dao"
	"gorm.io/gorm"
)

type FileInfoModel interface {
	InsertFileInfo(record *dao.FileInfo) error
}

type fileInfoModelImpl struct {
	db *gorm.DB
}

func NewFileInfoModel(db *gorm.DB) FileInfoModel {
	return &fileInfoModelImpl{db: db}
}

func (m *fileInfoModelImpl) InsertFileInfo(record *dao.FileInfo) error {
	tx := m.db
	return tx.Create(record).Error
}
