package model

import (
	"gitlab.senseauto.com/apcloud/library/common-go/client/sqldb"
)

func AutoMigrate(client *sqldb.Client, tables []interface{}) {
	if err := client.AutoMigrate(tables...); err != nil {
		panic(err)
	}
}
