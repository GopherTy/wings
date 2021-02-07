package cmd

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/gopherty/wings/common"
	"github.com/gopherty/wings/common/conf"
	"github.com/gopherty/wings/common/db"
	"github.com/gopherty/wings/common/logger"
	"github.com/gopherty/wings/module"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run wings",
	Run: func(cmd *cobra.Command, args []string) {
		// init common module
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		registers := []common.IRegister{
			conf.Register{},
			logger.Register{},
			db.Register{},
		}
		for _, reg := range registers {
			if err := reg.Regist(); err != nil {
				log.Fatalf("%s register failed. %s %v %s\n", reg.Name(), red, err, reset)
			}
		}

		// run grpc gateway server
		if err := runGateway(); err != nil {
			logger.Instance().Sugar().Fatalf("server serve failed. %s %v %s\n", red, err, reset)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runServer(ctx context.Context, gw *runtime.ServeMux) (err error) {
	// run grpc serve
	l, err := net.Listen("tcp", conf.Instance().Server.Address)
	if err != nil {
		logger.Instance().Sugar().Fatalf("server listen failed. %s %v %s \n", red, err, reset)
	}
	srv := grpc.NewServer()
	if err = module.InitManager(ctx, gw, srv); err != nil {
		logger.Instance().Sugar().Fatalf("init grpc module manager failed. %s %v  %s\n", red, err, reset)
	}
	logger.Instance().Sugar().Infof("grpc server running in %s", conf.Instance().Server.Address)
	if err := srv.Serve(l); err != nil {
		logger.Instance().Sugar().Fatalf("server serve failed. %s %v %s\n", red, err, reset)
	}
	return
}

func runGateway() (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	go runServer(ctx, mux)

	logger.Instance().Sugar().Infof("grpc gateway server running in %s", conf.Instance().GatewayServer.Address)
	return http.ListenAndServe(conf.Instance().GatewayServer.Address, mux)
}
