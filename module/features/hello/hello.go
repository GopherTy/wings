package hello

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/gopherty/wings/common/conf"
	pb "github.com/gopherty/wings/pb/hello"
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
func (Hello) RegisterServer(srv *grpc.Server) {
	pb.RegisterHelloServiceServer(srv, &helloService{})
}

func (Hello) RegisterHandler(mux *runtime.ServeMux) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	return pb.RegisterHelloServiceHandlerFromEndpoint(context.Background(), mux, conf.Instance().Server.Address, opts)
}
