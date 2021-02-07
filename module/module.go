package module

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// IModule 模块接口
type IModule interface {
	Name() string
	Init() error
	Registry(context.Context, *runtime.ServeMux, *grpc.Server) error
}
