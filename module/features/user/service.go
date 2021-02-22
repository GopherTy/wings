package user

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/gopherty/wings/common/db"
	pb "github.com/gopherty/wings/pb/user"
	"github.com/gopherty/wings/pkg/token"
)

type userService struct {
	pb.UnimplementedUserServiceServer
}

// Login  authorization
func (userService) Login(ctx context.Context, req *pb.LoginRequest) (resp *pb.LoginResponse, err error) {
	if req.GetUser() == "" || req.GetPassword() == "" {
		err = status.Error(codes.InvalidArgument, "user or password can not be null")
		return
	}

	// valid administrator user and password
	u := &db.Administrator{
		User:     req.GetUser(),
		Password: req.GetPassword(),
	}
	ok, err := db.Engine().Get(u)
	if err != nil {
		return
	}
	if !ok {
		err = status.Error(codes.Unauthenticated, "user or password not correct")
		return
	}

	var addr string
	at := time.Now()
	p, ok := peer.FromContext(ctx)
	if ok {
		addr = p.Addr.String()
	} else {
		addr = "unknown"
	}
	_, err = db.Engine().Cols(db.AdministratorColIPAddress, db.AdministratorColLastLogin).ID(u.ID).Update(&db.Administrator{
		IPAddress: addr,
		LastLogin: at,
	})
	if err != nil {
		return
	}

	t, err := token.New(u.ID)
	if err != nil {
		return
	}
	resp = &pb.LoginResponse{
		AccessToken:  t.Access,
		RefreshToken: t.Refresh,
		LastIp:       u.IPAddress,
		LastTime:     u.LastLogin.Local().String(),
	}
	return
}
