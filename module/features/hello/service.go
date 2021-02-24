package hello

import (
	"context"

	pb "github.com/gopherty/wings/pb/hello"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type helloService struct {
	pb.UnimplementedHelloServiceServer
}

func (helloService) SayHello(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloResponse, err error) {
	if req.GetName() == "" {
		err = status.Errorf(codes.InvalidArgument, "name cannot be null")
		return
	}
	return &pb.HelloResponse{Message: "Hello " + req.GetName()}, nil
}
