package constant

import "time"

const (
	AppName          = "app.datacollector"
	AppCollectorPath = "/noauth/v1/collector"
	FlagConfig       = "config"
	FlagVerbose      = "verbose"

	KeyRedisRetryJobs  = AppName + ".retryjobs"
	DurationExpire     = 24 * time.Hour
	MaxRequestTimeout  = 30 * time.Second
	MaxInsertCount     = 1000
	DefaultHTTPTimeout = 30

	CfgHTTPPort        = "http_port"
	CfgRPCPort         = "rpc_port"
	CfgSQLDB           = "sqldb"
	CfgRedis           = "redis"
	CfgRedisType       = "single"
	CfgRedisDB         = 9
	CfgOSS             = "oss"
	CfgDCP             = "DCP"
	CfgCloud           = "cloud"
	CfgMonitor         = "monitor"
	CfgAlert           = "alert"
	CfgInfraMonth      = "InfraMonth"
	CfgDCPDailyStorage = "DCPDailyStorage"
	CfgPrice           = "Price"
	CfgAOSSStorage     = "aoss"
	CfgSimulation      = "simulation"
	CfgOnes            = "ones"

	CfgHTTPAlert              = "http.alert"
	CfgHTTPAlertKey           = "alert"
	CfgHTTPCloudManagementKey = "management_url"
	CfgHTTPCloudAFSKey        = "afs_url"
	CfgHTTPCloudACSKey        = "acs_url"
	CfgHTTPCloudAOSSKey       = "aoss_url"
	CfgHTTPCloudIAMKey        = "iam_url"
	CfgHTTPDCP                = "http.DCP"
	CfgHTTPCloud              = "http.cloud"
	CfgHTTPDataProMeta        = "http.datapro-meta"
	CfgGRPCFieldTest          = "grpc.fieldtest"
	CfgHTTPOnes               = "http.ones"
	CfgSyncDBTruthValue       = "syncdb.truthvalue"
)

const (
	I18nEnFilePath = "locales/i18n.en.toml"
	I18nZhFilePath = "locales/i18n.zh.toml"
)

const (
	ClusterTypeNone  = -1
	ClusterTypeDcp   = 0
	ClusterTypeCloud = 1

	KeyURL = "url"

	ServiceDcpExporter = "dcpExporter"

	IDCollectorGpu            = "GPUCollector"
	IDCollectorGpuRange       = "GPURangeCollector"
	IDCollectorDcpStorage     = "DCPStorageDailyCollector"
	IDCollectorCloudStorage   = "CloudStorageCollector"
	IDCollectorInfraMonthBill = "InfraMonthBillCollector"
	IDCollectorSimulation     = "SimulationCollector"
	IDCollectorFieldTest      = "FieldTestCollector"
	IDCollectorOnes           = "OnesCollector"
	IDCollectorTruthValue     = "TruthValueCollector"

	FmtAlertPlus = ">%s:<font color=\"comment\">%s</font>"

	AlertServiceDcpDailyStorage = "DCP集群存储日账单导入"
	AlertServiceInfraMonthBill  = "基础设施月账单导入"
	AlertServicePrice           = "基础设施单价导入"
	DeviceTypeGpu               = 0
	DeviceTypeStorage           = 1
	DeviceTypeCPU               = 2
)

const (
	DcpAllPath = "/dcpgpuall"

	MethodGet  = "GET"
	MethodPost = "POST"

	FmtInfraMonth         = "【账单沙盘-V13\\.0】大装置资源账单管理_%s.xlsx"
	FmtIagStorage         = "cloud_storage_using_report_%s.xlsx"
	FmtInfraResourcePrice = "基础设置资源单价_%s.xlsx"

	DcpDailyStorageFile = 0
	InfraMonthFile      = 1
)

type StorageType int

const (
	TypeAFS StorageType = iota
	TypeAOSS
	TypeACS
)

const (
	RetryMsgBatchRetch = 3
)

const (
	HTTPCodeSucess           = "success.ok"
	HTTPCodeSuccessOnes      = 0
	HTTPCodeSuccessCollector = ""
)

const (
	TypeGpuRangeDaily int = iota
	TypeGpuRangeWeekly
	TypeGpuRangeMonthly

	FullTimeTemplate = "2006-01-02 15:04:05"
	DateTemplate     = "2006-01-02"
	MonthTemplate    = "2006-01"
)

type WorkflowJobStatus string

const (
	JobPending    WorkflowJobStatus = "Pending"
	JobRunning    WorkflowJobStatus = "Running"
	JobPause      WorkflowJobStatus = "Pause"
	JobFinished   WorkflowJobStatus = "Finished"
	JobUnfinished WorkflowJobStatus = "Unfinished"
)

type OnesIssueType string

const (
	IssueTypeRequirement OnesIssueType = "需求"
	IssueTypeTask        OnesIssueType = "任务"
	IssueTypeDefect      OnesIssueType = "缺陷"
	IssueTypeQuestion    OnesIssueType = "question"
)

type OnesTaskStatus string

const (
	TaskStatusRejected   OnesTaskStatus = "Rejected"
	TaskStatusDone       OnesTaskStatus = "Done"
	TaskStatusCancelled  OnesTaskStatus = "Cancelled"
	TaskStatusAccepted   OnesTaskStatus = "Accepted"
	TaskStatusTesting    OnesTaskStatus = "Testing"
	TaskStatusInReview   OnesTaskStatus = "In Review"
	TaskStatusInProgress OnesTaskStatus = "In Progress"
	TaskStatusToDo       OnesTaskStatus = "To Do"
	TaskStatusPending    OnesTaskStatus = "Pending"
	TaskStatusOpen       OnesTaskStatus = "Open"
	TaskStatusReview     OnesTaskStatus = "Review"
)

type OnesTaskPriority string

const (
	TaskPriorityLowest  OnesTaskPriority = "最低"
	TaskPriorityLow     OnesTaskPriority = "较低"
	TaskPriorityCommon  OnesTaskPriority = "普通"
	TaskPriorityHigh    OnesTaskPriority = "较高"
	TaskPriorityHighest OnesTaskPriority = "最高"

	TaskPriorityBlocker  OnesTaskPriority = "Blocker"
	TaskPriorityCritical OnesTaskPriority = "Critical"
	TaskPriorityMajor    OnesTaskPriority = "Major"
	TaskPriorityNormal   OnesTaskPriority = "Normal"
	TaskPriorityTrivial  OnesTaskPriority = "Trivial"
)
