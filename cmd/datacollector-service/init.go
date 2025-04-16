package main

import (
	"context"

	"github.com/spf13/cobra"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/lib/constant"
	cmclient "gitlab.senseauto.com/apcloud/library/common-go/client"
	cmsql "gitlab.senseauto.com/apcloud/library/common-go/client/sqldb"
	cmconfig "gitlab.senseauto.com/apcloud/library/common-go/config"
	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
)

func attachInitCommand(rootCmd *cobra.Command) {
	initCmd := &cobra.Command{
		Use: "init",
		Run: func(cmd *cobra.Command, args []string) {
			logger := cmlog.New(
				cmlog.WithAppName(constant.AppName),
			)
			ctx := cmlog.Inject(context.Background(), logger)

			cfgFile, err := cmd.Flags().GetString(constant.FlagConfig)
			if err != nil {
				logger.Panic(err)
			}

			if err := cmconfig.Load(cfgFile); err != nil {
				logger.Panic(err)
			}

			initSQLDB(ctx, cmd)
		},
	}
	initCmd.Flags().String(constant.FlagConfig, "conf/config.yaml", "set the path of configuration file")
	rootCmd.AddCommand(initCmd)
}

func initSQLDB(ctx context.Context, cmd *cobra.Command) {
	logger := cmlog.Extract(ctx)
	sqldbCfg := cmsql.Config{}

	if err := cmconfig.Global().UnmarshalKey(constant.CfgSQLDB, &sqldbCfg); err != nil {
		logger.Panic(err)
	}
	cmclient.SQLDB.Global(ctx, sqldbCfg)

	RegisterService(ctx)
	AutoMigrate(ctx, cmclient.SQLDB)
}
