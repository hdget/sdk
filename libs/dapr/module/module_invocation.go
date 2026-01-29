package module

import (
	reflectUtils "github.com/hdget/utils/reflect"
	"github.com/pkg/errors"
)

type InvocationModule interface {
	Module
	RegisterHandlers(functions map[string]InvocationFunction) error // 注册Handlers
	GetHandlers() []invocationHandler                               // 获取handlers
}

type invocationModuleImpl struct {
	Module
	self     any // 实际module实例
	handlers []invocationHandler
}

var (
	_ InvocationModule = (*invocationModuleImpl)(nil)
)

func NewInvocationModule(moduleObject any, app string, alias2handler map[string]InvocationFunction) error {
	// 首先实例化module
	module, err := asInvocationModule(moduleObject, app)
	if err != nil {
		return err
	}

	// 然后注册handlers
	err = module.RegisterHandlers(alias2handler)
	if err != nil {
		return err
	}

	// 最后注册module
	register(module)

	return nil
}

// RegisterHandlers 参数handlers为alias=>receiver.fnName, 保存为handler.id=>*invocationHandler
func (impl *invocationModuleImpl) RegisterHandlers(functions map[string]InvocationFunction) error {
	impl.handlers = make([]invocationHandler, 0)
	for handlerAlias, fn := range functions {
		impl.handlers = append(impl.handlers, impl.newInvocationHandler(impl.Module, handlerAlias, fn))
	}
	return nil
}

func (impl *invocationModuleImpl) GetKind() ModuleKind {
	return ModuleKindInvocation
}

func (impl *invocationModuleImpl) GetHandlers() []invocationHandler {
	return impl.handlers
}

// asInvocationModule 将一个any类型的结构体转换成InvocationModule
//
// e,g:
//
//		type v1_test struct {
//		  InvocationModule
//		}
//
//		 v := &v1_test{}
//		 im, err := asInvocationModule("app",v)
//	     if err != nil {
//	      ...
//	     }
//	     im.DiscoverHandlers()
func asInvocationModule(moduleObject any, app string) (InvocationModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &invocationModuleImpl{
		Module: m,
		self:   moduleObject,
	}

	// 初始化module
	err = reflectUtils.StructSet(moduleObject, (*InvocationModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "devops module: %+v", m)
	}

	module, ok := moduleObject.(InvocationModule)
	if !ok {
		return nil, errors.New("invalid invocation module")
	}

	return module, nil
}

func (impl *invocationModuleImpl) newInvocationHandler(module Module, handlerAlias string, fn InvocationFunction) invocationHandler {
	return &invocationHandlerImpl{
		handlerAlias: handlerAlias,
		handlerName:  reflectUtils.GetFuncName(fn),
		module:       module,
		fn:           fn,
	}
}
