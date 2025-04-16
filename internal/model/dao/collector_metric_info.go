package dao

type CollectorMetricInfo struct {
	ID          uint64  `json:"id" gorm:"primaryKey;type:bigint;comment:编号"`
	CollectorID string  `json:"collector_id" gorm:"type:varchar(50);index:collector_sync_idx_collector_id;comment:采集服务编码"`
	Method      string  `json:"method" gorm:"type:varchar(50);comment:collector方法"`
	Duration    float64 `json:"duration" gorm:"type:bigint;comment:支持时间"`
	RuntimeNum  int32   `json:"runtime_num" gorm:"type:int;comment:当前运行的tick或者cron个数"`
	CreateAt    int64   `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;comment:创建时间"`
	Status      int32   `json:"status" gorm:"type:int;comment:状态，-1、undone 0、success,1、error,2、timeout"`
	ExtraInfo   string  `json:"extra_info" gorm:"type:text;comment:额外信息"`
}

// TableName sets the insert table name for this struct type
func (v *CollectorMetricInfo) TableName() string {
	return "p_collector_metric_info"
}
