package cmd

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/gopherty/wings/common"
	"github.com/gopherty/wings/common/conf"
	"github.com/gopherty/wings/common/db"
	"github.com/gopherty/wings/common/logger"
	"github.com/gopherty/wings/module"
	"github.com/gopherty/wings/pkg/colors"
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

		addr := conf.Instance().Server.Address
		l, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Instance().Sugar().Fatalf("failed to listen %s. %v", addr, err)
		}
		srv := newServer()

		logger.Instance().Sugar().Infof("server work on %s", addr)
		srv.serve(l)

		// start service
		// if err := serve(); err != nil {
		// 	logger.Instance().Sugar().Fatalf("server serve failed. %s %v %s\n", red, err, reset)
		// }
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

	logger.Instance().Sugar().Infof("grpc gateway server running in %s", conf.Instance().Gateway.Address)
	return http.ListenAndServe(conf.Instance().Gateway.Address, mux)
}

// wings server
type server struct {
	s *grpc.Server

	gateway *runtime.ServeMux // grpc gateway, it reads a gRPC service definition and generates a reverse-proxy server which translates a RESTful JSON API into gRPC
}

func (srv *server) serve(l net.Listener) error {
	if srv.s == nil {
		panic("grpc server is nil")
	}

	// translates a RESTful JSON API into gRPC
	if srv.gateway != nil {
		go func() {
			cnf := conf.Instance().Gateway
			var err error
			defer logger.Instance().Sugar().Errorf("failed to listen and serve. {%v}", err)

			if cnf.H2 {
				if cnf.CertFile != "" && cnf.KeyFile != "" {
					logger.Instance().Sugar().Infof("h2 work on %s", cnf.Address)

					err = http.ListenAndServeTLS(cnf.Address, cnf.CertFile, cnf.KeyFile, srv.gateway)
				} else {
					logger.Instance().Sugar().Infof("h2c work on %s", cnf.Address)

					handler := h2c.NewHandler(srv.gateway, &http2.Server{})
					err = http.ListenAndServe(cnf.Address, handler)
				}
			} else {
				logger.Instance().Sugar().Infof("http work on %s", cnf.Address)

				err = http.ListenAndServe(cnf.Address, srv.gateway)
			}
		}()
	}

	return srv.s.Serve(l)
}

func newServer() *server {
	srv := new(server)

	// server dev mode add request logs interceptor
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logsUnaryInterceptor),
		grpc.StreamInterceptor(logsStreamInterceptor),
	}
	srv.s = grpc.NewServer(opts...)

	if conf.Instance().Gateway != nil {
		var opts []runtime.ServeMuxOption
		srv.gateway = runtime.NewServeMux(opts...)
	}

	return srv
}

func logsUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	at := time.Now()
	var addr string
	if p, ok := peer.FromContext(ctx); ok {
		addr = p.Addr.String()
	} else {
		addr = "TODO"
	}
	resp, err = handler(ctx, req)
	if err != nil {
		if conf.Instance().Logger.LogsPath != "" { // logs file
			logger.Instance().Sugar().Errorf(" %s | %s | %v  {%v}", addr, info.FullMethod, time.Since(at), err)
		} else { // stdout
			logger.Instance().Sugar().Errorf(" %s | %s %s %s | %v  {%v}", addr, colors.Red, info.FullMethod, colors.Reset, time.Since(at), err)
		}
		return
	}

	if conf.Instance().Logger.LogsPath != "" { // logs file
		logger.Instance().Sugar().Infof(" %s | %s | %v ", addr, info.FullMethod, time.Since(at))
	} else { // stdout
		logger.Instance().Sugar().Infof(" %s | %s %s %s | %v ", addr, colors.Green, info.FullMethod, colors.Reset, time.Since(at))
	}
	return
}

func logsStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	at := time.Now()
	var addr string
	ctx := ss.Context()
	if p, ok := peer.FromContext(ctx); ok {
		addr = p.Addr.String()
	} else {
		addr = "TODO"
	}
	err := handler(srv, ss)
	if err != nil {
		if conf.Instance().Logger.LogsPath != "" { // logs file
			logger.Instance().Sugar().Errorf(" %s | %s | %v  {%v}", addr, info.FullMethod, time.Since(at), err)
		} else { // stdout
			logger.Instance().Sugar().Errorf(" %s | %s %s %s | %v  {%v}", addr, colors.Red, info.FullMethod, colors.Reset, time.Since(at), err)
		}
		return err
	}

	if conf.Instance().Logger.LogsPath != "" { // logs file
		logger.Instance().Sugar().Infof(" %s | %s | %v ", addr, info.FullMethod, time.Since(at))
	} else { // stdout
		logger.Instance().Sugar().Infof(" %s | %s %s %s | %v ", addr, colors.Green, info.FullMethod, colors.Reset, time.Since(at))
	}
	return nil
}
