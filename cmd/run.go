package cmd

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

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

		// start service
		if err := serve(); err != nil {
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
		return
	}

	// add
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		at := time.Now()
		var addr string
		if p, ok := peer.FromContext(ctx); ok {
			addr = p.Addr.String()
		} else {
			addr = "unknown address"
		}
		resp, err = handler(ctx, req)
		if err == nil {
			logger.Instance().Sugar().Infof(" %s | %s %s %s | %v ", addr, "\033[97;42m", info.FullMethod, "\033[0m", time.Since(at))
		} else {
			logger.Instance().Sugar().Errorf(" %s | %s %s %s | %v  {%v}", addr, "\033[97;41m", info.FullMethod, "\033[0m", time.Since(at), err)
		}
		return
	}))
	opts = append(opts, grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		at := time.Now()
		var addr string
		ctx := ss.Context()
		if p, ok := peer.FromContext(ctx); ok {
			addr = p.Addr.String()
		} else {
			addr = "unknown address"
		}
		err = handler(srv, ss)
		if err == nil {
			logger.Instance().Sugar().Infof(" %s | %s %s %s | %v ", addr, "\033[97;42m", info.FullMethod, "\033[0m", time.Since(at))
		} else {
			logger.Instance().Sugar().Errorf(" %s | %s %s %s | %v  {%v}", addr, "\033[97;41m", info.FullMethod, "\033[0m", time.Since(at), err)
		}
		return err
	}))
	srv := grpc.NewServer(opts...)
	if err = module.Init(ctx, gw, srv); err != nil {
		logger.Instance().Sugar().Fatalf("init grpc module manager failed. %s %v  %s\n", red, err, reset)
		return
	}

	logger.Instance().Sugar().Infof("grpc server running in %s", conf.Instance().Server.Address)
	if err := srv.Serve(l); err != nil {
		logger.Instance().Sugar().Fatalf("grpc server serve failed. %s %v %s\n", red, err, reset)
	}
	return
}

func serve() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	go runServer(ctx, mux)

	logger.Instance().Sugar().Infof("grpc gateway server running in %s", conf.Instance().GatewayServer.Address)
	return http.ListenAndServe(conf.Instance().GatewayServer.Address, mux)
}
