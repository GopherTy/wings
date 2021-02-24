package hello

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/gopherty/wings/common/conf"
	pb "github.com/gopherty/wings/pb/hello"
	"github.com/gopherty/wings/pkg/token"
)

// Hello .
type Hello struct {
}

// Name .
func (Hello) Name() string {
	return "Features.Hello"
}

// Init .
func (Hello) Init() error {
	return nil
}

// Registry registry hello module
func (Hello) Registry(ctx context.Context, mux *runtime.ServeMux, srv *grpc.Server) error {
	pb.RegisterHelloServiceServer(srv, &helloService{})

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			err := token.ValidToken(ctx, token.AccessKeyFunc)
			if err != nil {
				return err
			}
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	}

	return pb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, conf.Instance().Server.Address, opts)
}
