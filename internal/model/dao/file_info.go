package dao

type FileInfo struct {
	ID       uint64 `json:"id" gorm:"primaryKey;type:bigint;comment:编号"`
	Name     string `json:"file_name" gorm:"type:varchar(50);index:file_info_idx_name;comment:导入文件的文件名"`
	Type     int32  `json:"type" gorm:"type:int;index:file_info_idx_type;comment:导入文件的类型: 0: dcp 1: month 2: price"`
	URL      string `json:"url" gorm:"type:varchar(200);comment:导入文件的链接地址"`
	Date     int64  `json:"file_date" gorm:"type:bigint;comment:文件数据时间"`
	MD5      string `json:"md5" gorm:"type:varchar(50);comment:导入文件的MD5值"`
	CreateAt int64  `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;comment:创建时间"`
}

// TableName sets the insert table name for this struct type
func (v *FileInfo) TableName() string {
	return "c_file_info"
}
