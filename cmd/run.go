package cmd

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	"github.com/gopherty/wings/common"
	"github.com/gopherty/wings/common/conf"
	"github.com/gopherty/wings/common/db"
	"github.com/gopherty/wings/common/logger"
	"github.com/gopherty/wings/module"
	"github.com/gopherty/wings/pkg/colors"
	"github.com/gopherty/wings/pkg/token"
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
				log.Fatalf("%s register failed. %s %v %s\n", reg.Name(), colors.Red, err, colors.Reset)
			}
		}

		addr := conf.Instance().Server.Address
		l, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Instance().Sugar().Fatalf("failed to listen %s. %v", addr, err)
		}
		srv := newServer()

		err = module.Init(srv.s, srv.gateway)
		if err != nil {
			logger.Instance().Sugar().Fatalf("init moudle failed %v", err)
		}

		logger.Instance().Sugar().Infof("server work on %s", addr)
		if err = srv.serve(l); err != nil {
			logger.Instance().Sugar().Fatalf("serve failed. %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
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

			if cnf.H2 {
				if cnf.CertFile != "" && cnf.KeyFile != "" {
					logger.Instance().Sugar().Infof("h2 work on %s", cnf.Address)

					if err = http.ListenAndServeTLS(cnf.Address, cnf.CertFile, cnf.KeyFile, srv.gateway); err != nil {
						logger.Instance().Sugar().Errorf("failed to listen and serve. {%v}", err)
					}
				} else {
					logger.Instance().Sugar().Infof("h2c work on %s", cnf.Address)

					handler := h2c.NewHandler(srv.gateway, &http2.Server{})
					if err = http.ListenAndServe(cnf.Address, handler); err != nil {
						logger.Instance().Sugar().Errorf("failed to listen and serve. {%v}", err)
					}
				}
			} else {
				logger.Instance().Sugar().Infof("http work on %s", cnf.Address)

				if err = http.ListenAndServe(cnf.Address, srv.gateway); err != nil {
					logger.Instance().Sugar().Errorf("failed to listen and serve. {%v}", err)
				}
			}
		}()
	}

	return srv.s.Serve(l)
}

func newServer() *server {
	srv := new(server)

	// server dev mode add request logs interceptor
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	}

	if conf.Instance().Server.CertFile != "" && conf.Instance().Server.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(conf.Instance().Server.CertFile, conf.Instance().Server.KeyFile)
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	}

	srv.s = grpc.NewServer(opts...)

	if conf.Instance().Gateway != nil {
		var opts []runtime.ServeMuxOption
		srv.gateway = runtime.NewServeMux(opts...)
	}

	return srv
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	at := time.Now()
	var addr string
	if p, ok := peer.FromContext(ctx); ok {
		addr = p.Addr.String()
	} else {
		addr = "TODO"
	}

	if info.FullMethod != "/user.UserService/Login" {
		err = token.ValidToken(ctx, token.AccessKeyFunc)
		if err != nil {
			if conf.Instance().Logger.LogsPath != "" { // logs file
				logger.Instance().Sugar().Errorf(" %s | %s | %v  {%v}", addr, info.FullMethod, time.Since(at), err)
			} else { // stdout
				logger.Instance().Sugar().Errorf(" %s | %s %s %s | %v  {%v}", addr, colors.Red, info.FullMethod, colors.Reset, time.Since(at), err)
			}
			return
		}
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

func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	at := time.Now()
	var addr string
	ctx := ss.Context()
	if p, ok := peer.FromContext(ctx); ok {
		addr = p.Addr.String()
	} else {
		addr = "TODO"
	}

	err := token.ValidToken(ctx, token.AccessKeyFunc)
	if err != nil {
		if conf.Instance().Logger.LogsPath != "" { // logs file
			logger.Instance().Sugar().Errorf(" %s | %s | %v  {%v}", addr, info.FullMethod, time.Since(at), err)
		} else { // stdout
			logger.Instance().Sugar().Errorf(" %s | %s %s %s | %v  {%v}", addr, colors.Red, info.FullMethod, colors.Reset, time.Since(at), err)
		}
		return err
	}
	err = handler(srv, ss)
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
