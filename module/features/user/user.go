package user

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/gopherty/wings/module/pb/user"
)

// User grpc feature user
type User struct {
	pb.UnimplementedUserServer
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
func (u User) Registry(srv *grpc.Server) {
	pb.RegisterUserServer(srv, u)
}

// Login .
func (User) Login(ctx context.Context, req *pb.LoginRequest) (resp *pb.LoginResponse, err error) {
	return
}
