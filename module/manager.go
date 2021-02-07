package module

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/gopherty/wings/common/logger"
	"github.com/gopherty/wings/module/features/user"
)

// Manager 模块控制器
type Manager struct {
	srv     *grpc.Server
	modules map[string]IModule
}

func (m *Manager) reset(s *grpc.Server) {
	m.srv = s

	m.modules = map[string]IModule{
		"User": user.User{},
	}
}

// default module manager
var m *Manager

// some default error
var (
	ErrServerNotAllowed = errors.New("grpc server is nil")
)

// InitManager initialization module manager
func InitManager(ctx context.Context, gw *runtime.ServeMux, srv *grpc.Server) (err error) {
	if srv == nil {
		return ErrServerNotAllowed
	}
	m = new(Manager)
	m.reset(srv)

	for _, v := range m.modules {
		//  internal module init
		if err = v.Init(); err != nil {
			logger.Instance().Sugar().Errorf(" %s init failed. %v", v.Name(), err)
			continue
		}
		v.Registry(ctx, gw, srv)
	}
	return
}
