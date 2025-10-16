package module

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	reflectUtils "github.com/hdget/utils/reflect"
	"github.com/pkg/errors"
)

type HealthModule interface {
	Module
	GetHandler() common.HealthCheckHandler
}

type healthModuleImpl struct {
	Module
	fn HealthCheckFunction
}

type HealthCheckFunction func(context.Context) error

var (
	_                        HealthModule = (*healthModuleImpl)(nil)
	EmptyHealthCheckFunction              = func(ctx context.Context) (err error) { return nil }
)

// NewHealthModule 注册健康模块
func NewHealthModule(moduleObject any, app string, fn HealthCheckFunction) error {
	// 首先实例化module
	module, err := asHealthModule(moduleObject, app, fn)
	if err != nil {
		return err
	}

	// 最后注册module
	register(module)

	return nil
}

// asHealthModule 将一个any类型的结构体转换成HealthModule
//
// e,g:
//
//		type v1_test struct {
//		  HealthModule
//		}
//
//		 v := &v1_test{}
//		 im, err := asHealthModule("app",v)
//	     if err != nil {
//	      ...
//	     }
func asHealthModule(moduleObject any, app string, fn HealthCheckFunction) (HealthModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &healthModuleImpl{
		Module: m,
		fn:     fn,
	}

	// 初始化module
	err = reflectUtils.StructSet(moduleObject, (*HealthModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "devops module: %+v", m)
	}

	module, ok := moduleObject.(HealthModule)
	if !ok {
		return nil, errors.New("invalid health module")
	}

	return module, nil
}

func (impl *healthModuleImpl) GetKind() ModuleKind {
	return ModuleKindHealth
}

func (impl *healthModuleImpl) GetHandler() common.HealthCheckHandler {
	if impl.fn == nil {
		return EmptyHealthCheckFunction
	}

	return func(ctx context.Context) error {
		return impl.fn(ctx)
	}
}
