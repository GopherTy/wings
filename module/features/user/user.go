package user

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/gopherty/wings/common/conf"
	pb "github.com/gopherty/wings/module/pb/user"
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

// Registry regist user module
func (User) Registry(ctx context.Context, mux *runtime.ServeMux, srv *grpc.Server) error {
	pb.RegisterUserServiceServer(srv, userService{})

	opts := []grpc.DialOption{grpc.WithInsecure()}
	return pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, conf.Instance().Server.Address, opts)
}
