package module

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// IModule 模块接口
type IModule interface {
	Init() error
	Name() string
	RegisterServer(*grpc.Server)             // grpc server register
	RegisterHandler(*runtime.ServeMux) error // grpc gateway handler register
}
