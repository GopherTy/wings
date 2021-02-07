package user

import (
	"context"

	pb "github.com/gopherty/wings/module/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login .
func (User) Login(ctx context.Context, req *pb.LoginRequest) (resp *pb.LoginResponse, err error) {
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
