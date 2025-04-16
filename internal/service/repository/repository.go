package repository

import (
	"context"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/lib/constant"
	cmsql "gitlab.senseauto.com/apcloud/library/common-go/client/sqldb"
	cmconfig "gitlab.senseauto.com/apcloud/library/common-go/config"
	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
)

const (
	dbName = "bi_service"
	dnEnv  = "MYSQL_DATABASE_BISERVICE"
)

type holder struct {
	common ComRepository
	core   CoreRepository
}

var instance holder = holder{
	common: ComRepository{},
	core:   CoreRepository{},
}

func createDBClient(ctx context.Context) (*cmsql.Client, error) {
	logger := cmlog.Extract(ctx)
	sqldbCfg := cmsql.Config{
		DatabaseEnv: dnEnv,
	}

	if err := cmconfig.Global().UnmarshalKey(constant.CfgSQLDB, &sqldbCfg); err != nil {
		logger.Panic(err)
	}
	sqldbCfg.DBName = dbName
	clent := cmsql.New(ctx, sqldbCfg)
	return clent, nil
}

func AutoMigrate(ctx context.Context) {
	client, err := createDBClient(ctx)
	if err == nil {
		tables := make([]interface{}, 0)
		tables = append(tables, instance.common.getMigrateTables()...)
		tables = append(tables, instance.core.getMigrateTables()...)
		if err := client.AutoMigrate(tables...); err != nil {
			panic(err)
		}
	}
}

func CreateRepository(ctx context.Context) (*ComRepository, *CoreRepository, error) {
	client, err := createDBClient(ctx)
	if err == nil {
		err = instance.common.init(ctx, client)
		if err == nil {
			err = instance.core.init(ctx, client)
		}
	}
	return &instance.common, &instance.core, err
}
