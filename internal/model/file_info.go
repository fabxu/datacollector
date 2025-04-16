package model

import (
	"errors"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	auth_filter "gitlab.senseauto.com/apcloud/library/common-go/auth/filter"
	"gorm.io/gorm"
)

type FileInfoModel interface {
	InsertFileInfo(record *dao.FileInfo) error
	GetLatestFileInfo(filter FileInfoFilter, columns []string) (*dao.FileInfo, error)
}

type fileInfoModelImpl struct {
	db *gorm.DB
}

type FileInfoFilter struct {
	ID       *auth_filter.UInt64Field
	Type     *auth_filter.Int32Field
	Name     *auth_filter.StringField
	CreateAt *auth_filter.UInt64Field
}

func NewFileInfoModel(db *gorm.DB) FileInfoModel {
	return &fileInfoModelImpl{db: db}
}

func (m *fileInfoModelImpl) InsertFileInfo(record *dao.FileInfo) error {
	tx := m.db
	return tx.Create(record).Error
}

func (m *fileInfoModelImpl) GetLatestFileInfo(filter FileInfoFilter, columns []string) (*dao.FileInfo, error) {
	tx := m.db.Model(&dao.FileInfo{})

	if len(columns) > 0 {
		// 如果columns不空，只返回指定列
		tx = tx.Select(columns)
	}

	tx = filter.ID.ToAndWhere(tx, "id")
	tx = filter.Name.ToAndWhere(tx, "name")
	tx = filter.Type.ToAndWhere(tx, "type")
	tx.Order("date DESC")

	fileInfo := dao.FileInfo{}

	result := tx.First(&fileInfo)
	if result.Error != nil {
		// 如果记录没有找到，返回空记录，不报错
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &fileInfo, nil
		}

		return nil, result.Error
	}

	return &fileInfo, nil
}
