package user

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/gopherty/wings/common/conf"
	gw "github.com/gopherty/wings/module/pb/user"
	pb "github.com/gopherty/wings/module/pb/user"
)

// User grpc feature user
type User struct {
	pb.UnimplementedUserServiceServer
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
func (u User) Registry(ctx context.Context, mux *runtime.ServeMux, srv *grpc.Server) (err error) {
	pb.RegisterUserServiceServer(srv, u)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	return gw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, conf.Instance().Server.Address, opts)
}
