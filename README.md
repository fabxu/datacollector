# datacollector

## 1、描述
bi看板用来统计、监控智驾各部门使用各型号gpu情况，包括使用率、年度财务结算等。后看板业务延伸，
统计了各部门使用存储（afs、aoss、acs）情况。再后来看板就统计了智驾云业务情况，例如路测、真值、甚至项目。
前端展示使用了帆软，由产品经理来配置报表。
后端需要起一个服务，使用定时任务的方式定时获取数据存库。

datacollector 是后端的框架，每增加一个采集项都需要注册到datacollector中，并实现工厂方法。
除此之外datacollector 还提供了采集预警、采集监控、心跳检测等功能。

## 2、方案设计
[技术方案]()


## 3、抽象方法
type Collector interface {
    GetMigrateTables() []interface{}
    Init(ctx *CollectorContext) error
    GetID() string
    Process(ctx context.Context, msgType int32, arg string) (interface{}, *monitor.MonitorMsg, error)
    Save(ctx context.Context, data interface{}) error
    GetTickOpts() (tick TickFunc, options []Option)
}

任何需要注册的collector都需要实现上述方法
### GetMigrateTables
获取需要需要自动迁移的表。后端采用的是gorm 框架数据库，表结构有改变的话都可以AutoMigrate 到线上
### Init
collector初始化。常见的是读取本地配置解析到golang备用
### GetID
获取collector id。例如gpu collector，storage collector。collector 的名字可以用来后续监控预警
### Process
collector 的采集逻辑处理。大部分数据都是通过调用api获取，例如专有云和公有云的数据，ones的数据，也有从es或者其他平台获取。
### Save
将上一步采集的数据存库。采集完的数据需要按各自要求存库，特定数据需要存es，高频采集的数据需要进行日平均等
### GetTickOpts
获取tick或者cron 参数，允许提供tick执行函数。基本上每个collector 都需要提供cron 参数，之前cron是用k8s cron组件做的，但是随着collector 越来越多，cron 变得不好管理，故由自身管理自己的cron.

## 4、更多参考文档