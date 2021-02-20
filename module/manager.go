package module

import (
	"context"
	"errors"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

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
	ctx context.Context

	interceptor Interceptor

	mu      sync.Mutex
	modules map[string]IModule
}

// Enable enable module
func (m *Manager) Enable(module IModule) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.modules == nil {
		m.modules = make(map[string]IModule)
	}

	name := module.Name()
	if _, ok := m.modules[name]; ok {
		err = ErrModuleExists
		return
	}

	err = module.Init()
	if err != nil {
		return
	}
	err = module.Registry(m.ctx, m.mux, m.srv)
	if err != nil {
		return
	}

	m.modules[name] = module
	return
}

func (m *Manager) reset(ctx context.Context, srv *grpc.Server, mux *runtime.ServeMux) {
	m.ctx = ctx
	m.srv = srv
	m.mux = mux
}

// Init initialization module manager
func Init(ctx context.Context, mux *runtime.ServeMux, srv *grpc.Server) error {
	if srv == nil || mux == nil {
		return ErrServerNotAllowed
	}
	m = new(Manager)
	m.reset(ctx, srv, mux)

	return m.Enable(user.User{})
}

// Instance retrun module  manager
func Instance() *Manager {
	return m
}
