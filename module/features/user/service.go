package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/gopherty/wings/module/pb/user"
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

	resp = &pb.LoginResponse{
		AccessToken:  req.GetUser(),
		RefreshToken: req.GetPassword(),
	}
	return
}
