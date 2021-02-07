package module

import "google.golang.org/grpc"

// IModule 模块接口
type IModule interface {
	Name() string
	Init() error
	Registry(*grpc.Server)
}
