package dao

// CollectorSyncInfo defines the model of a sync info
type CollectorSyncInfo struct {
	ID           uint64 `json:"id" gorm:"primaryKey;type:bigint;comment:编号"`
	CollectorID  string `json:"collector_id" gorm:"type:varchar(50);index:collector_sync_idx_collector_id;comment:采集服务编码"`
	LastSyncTime int64  `json:"last_sync_time" gorm:"type:bigint;autoUpdateTime:milli;comment:最近更新时间"`
	CreateAt     int64  `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;comment:创建时间"`
}

// TableName sets the insert table name for this struct type
func (v *CollectorSyncInfo) TableName() string {
	return "c_collector_sync_info"
}
