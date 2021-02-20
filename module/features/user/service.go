package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/gopherty/wings/pb/user"
	"github.com/gopherty/wings/pkg/token"
)

type userService struct {
	pb.UnimplementedUserServiceServer
}

// Login .
func (userService) Login(ctx context.Context, req *pb.LoginRequest) (resp *pb.LoginResponse, err error) {
	if req.GetUser() == "" || req.GetPassword() == "" {
		err = status.Error(codes.InvalidArgument, "user or password can not be null")
		return
	}

	t, err := token.New(1)
	if err != nil {
		return
	}
	resp = &pb.LoginResponse{
		AccessToken:  t.Access,
		RefreshToken: t.Refresh,
	}
	return
}
