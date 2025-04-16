package main

import (
	"context"

	"github.com/spf13/cobra"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/lib/constant"
	cmclient "gitlab.senseauto.com/apcloud/library/common-go/client"
	cmsql "gitlab.senseauto.com/apcloud/library/common-go/client/sqldb"
	cmconfig "gitlab.senseauto.com/apcloud/library/common-go/config"
	cmi18n "gitlab.senseauto.com/apcloud/library/common-go/lib/i18n"
	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
	cmserver "gitlab.senseauto.com/apcloud/library/common-go/server"
	datacollector_api "gitlab.senseauto.com/apcloud/library/proto/api/datacollector-service/v1"
)

func attachRunCommand(rootCmd *cobra.Command) {
	runCmd := &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			logger := cmlog.New(
				// 请填写自己的服务名称，便于日志查询
				cmlog.WithAppName(constant.AppName),
			)
			if verbose, err := cmd.Flags().GetCount(constant.FlagVerbose); err != nil {
				logger.Panic(err)
			} else {
				logger.SetLevel(cmlog.Level(0 - verbose))
			}

			cfgFile, err := cmd.Flags().GetString(constant.FlagConfig)
			if err != nil {
				logger.Panic(err)
			}

			if err := cmconfig.Load(cfgFile); err != nil {
				logger.Panic(err)
			}

			if err := cmconfig.Global().BindPFlags(cmd.Flags()); err != nil {
				logger.Panic(err)
			}

			ctx := cmlog.Inject(context.Background(), logger)
			setups(ctx)

			RegisterService(ctx)
			StartService(ctx)
			collectorService := GetServiceInstance()
			// historyImportService := service.NewHistoryImportService(ctx)
			opts := []cmserver.Option{
				cmserver.WithLogger(logger),
				cmserver.AddShutdown(),
			}

			cmserver.New(opts...).
				ChainRPC(cmconfig.Global().GetInt(constant.CfgRPCPort),
					cmserver.RPCRegister{Register: datacollector_api.RegisterCollectServiceServer, Server: collectorService},
					// cmserver.RPCRegister{Register: datacollector_api.RegisterHistoryImportServiceServer, Server: historyImportService},
				).
				ChainHTTP(cmconfig.Global().GetInt(constant.CfgHTTPPort),
					cmserver.HTTPRegister{Register: datacollector_api.RegisterCollectServiceHandlerFromEndpoint},
					// cmserver.HTTPRegister{Register: datacollector_api.RegisterHistoryImportServiceHandlerFromEndpoint},
				).
				Run()
		},
	}
	runCmd.Flags().String(constant.FlagConfig, "conf/config.yaml", "set the path of configuration file")
	runCmd.Flags().CountP(constant.FlagVerbose, "v", "print verbose info")
	rootCmd.AddCommand(runCmd)
}

func setups(ctx context.Context) {
	logger := cmlog.Extract(ctx)

	cmi18n.InitI18nByFiles([]string{constant.I18nEnFilePath, constant.I18nZhFilePath})

	// 初始化sqldb
	sqldbCfg := cmsql.Config{}
	if err := cmconfig.Global().UnmarshalKey(constant.CfgSQLDB, &sqldbCfg); err != nil {
		logger.Panic(err)
	}

	cmclient.SQLDB.Global(ctx, sqldbCfg)
}
