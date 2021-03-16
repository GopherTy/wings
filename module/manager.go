package module

import (
	"errors"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/gopherty/wings/module/features/hello"
	"github.com/gopherty/wings/module/features/user"
)

// module manager
var m *Manager

// some default error
var (
	ErrServerNotAllowed = errors.New("grpc server is nil")
	ErrModuleExists     = errors.New("module has exists")
)

// Manager 服务模块控制器
type Manager struct {
	srv *grpc.Server
	mux *runtime.ServeMux

	interceptor Interceptor

	mu      sync.Mutex
	modules map[string]IModule
}

func (m *Manager) enable(modules ...IModule) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.modules == nil {
		m.modules = make(map[string]IModule)
	}

	for _, module := range modules {
		name := module.Name()
		if _, ok := m.modules[name]; ok {
			err = ErrModuleExists
			return
		}

		err = module.Init()
		if err != nil {
			return
		}
		module.RegisterServer(m.srv)

		if m.mux != nil {
			err = module.RegisterHandler(m.mux)
			if err != nil {
				return
			}
		}

		m.modules[name] = module
	}

	return
}

// Init initialization module manager
func Init(srv *grpc.Server, mux *runtime.ServeMux) (err error) {
	if srv == nil {
		return ErrServerNotAllowed
	}
	m = new(Manager)
	m.srv = srv
	m.mux = mux

	// register grpc service server
	return m.enable(hello.Hello{}, user.User{})
}

// Instance retrun module  manager
func Instance() *Manager {
	return m
}
