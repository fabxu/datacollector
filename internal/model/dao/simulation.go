package dao

import "github.com/fabxu/datacollector-service/internal/lib/constant"

type SimulationRecord struct {
	ID                uint64  `json:"id" gorm:"primaryKey;type:bigint;comment:编号"`
	JobID             int32   `json:"job_id" gorm:"type:int;"`
	Name              string  `json:"name" gorm:"type:varchar(255);comment:仿真任务名"`
	SpaceID           int32   `json:"space_id" gorm:"type:int;"`
	SpaceType         string  `json:"space_type" gorm:"type:varchar(255);comment:业务类型"`
	SpaceTypeCategory string  `json:"space_type_category" gorm:"type:varchar(255);comment:业务类型大类"`
	TotalDistance     float32 `json:"total_distance" gorm:"type:float;comment:总里程"`
	Status            int32   `json:"status" gorm:"type:int;index:simu_idx_status;comment:状态，0、finished,1、unfinished,2、pause,3、running"`
	CreateTime        int64   `json:"create_time" gorm:"type:bigint;autoCreateTime:milli;comment:创建时间"`
	StartTime         int64   `json:"start_time" gorm:"type:bigint;autoCreateTime:milli;comment:开始时间"`
	EndTime           int64   `json:"end_time" gorm:"type:bigint;autoCreateTime:milli;index:simu_idx_end_time;comment:结束时间"`
	Creator           string  `json:"creator" gorm:"type:varchar(50);comment:创建人"`
	ExtraInfo         string  `json:"extra_info" gorm:"type:text;comment:额外信息"`
}

// TableName sets the insert table name for this struct type
func (d *SimulationRecord) TableName() string {
	return "simulation_record"
}

type SimRecordResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Count int32            `json:"count"`
		Jobs  []JobsProperties `json:"jobs"`
	} `json:"data"`
}

type JobsProperties struct {
	JobID          int32                      `json:"jobId"`
	Name           string                     `json:"name"`
	SpaceID        int32                      `json:"spaceId"`
	SpaceIdcn      string                     `json:"spaceIdcn"`
	DataOrigin     string                     `json:"dataOrigin"`
	TotalDistance  float32                    `json:"totalDistance"`
	Status         constant.WorkflowJobStatus `json:"status"`
	CreateTime     string                     `json:"createAt"`
	StartTime      string                     `json:"startTime"`
	EndTime        string                     `json:"endTime"`
	Creator        string                     `json:"creator"`
	SimulationType string                     `json:"simulationType"`
}
