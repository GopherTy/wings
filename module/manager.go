package module

import (
	"errors"
	"strings"

	"google.golang.org/grpc"
)

// Manager 模块控制器
type Manager struct {
	srv     *grpc.Server
	modules map[string]IModule
}

func (m *Manager) reset(s *grpc.Server) {
	m.srv = s
	m.modules = make(map[string]IModule)
}

// default module manager
var m *Manager

// some default error
var (
	ErrServerNotAllowed = errors.New("grpc server is nil")
)

// InitManager initialization module manager
func InitManager(s *grpc.Server) (err error) {
	if s == nil {
		return ErrServerNotAllowed
	}
	m.reset(s)

	for _, v := range m.modules {
		//  internal module init
		if err = v.Init(); err != nil {
		}
		if strings.HasPrefix(v.Name(), "Feature") {
			v.GRPCFeatureRegistry(s)
		}
	}
	return
}
