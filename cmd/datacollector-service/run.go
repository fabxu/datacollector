package main

import (
	"context"

	"github.com/fabxu/datacollector-service/internal/lib/constant"
	cmclient "github.com/fabxu/lib/client"
	cmsql "github.com/fabxu/lib/client/sqldb"
	cmconfig "github.com/fabxu/lib/config"
	cmlog "github.com/fabxu/log"
	cmserver "github.com/fabxu/server"
	"github.com/spf13/cobra"
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
	// 初始化sqldb
	sqldbCfg := cmsql.Config{}
	if err := cmconfig.Global().UnmarshalKey(constant.CfgSQLDB, &sqldbCfg); err != nil {
		logger.Panic(err)
	}

	cmclient.SQLDB.Global(ctx, sqldbCfg)
}
