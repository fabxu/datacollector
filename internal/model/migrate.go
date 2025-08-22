package model

import (
	"github.com/fabxu/lib/client/sqldb"
)

func AutoMigrate(client *sqldb.Client, tables []interface{}) {
	if err := client.AutoMigrate(tables...); err != nil {
		panic(err)
	}
}
