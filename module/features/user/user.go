package user

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/gopherty/wings/common/conf"
	pb "github.com/gopherty/wings/pb/user"
)

// User grpc feature user
type User struct {
}

// Name module name
func (User) Name() string {
	return "Features.User"
}

// Init .
func (User) Init() error {
	return nil
}

// RegisterServer register user module
func (User) RegisterServer(srv *grpc.Server) {
	pb.RegisterUserServiceServer(srv, userService{})
}

func (User) RegisterHandler(mux *runtime.ServeMux) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	return pb.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, conf.Instance().Server.Address, opts)
}
