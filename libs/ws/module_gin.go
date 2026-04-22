package lib_ws

import (
	"github.com/gin-gonic/gin"
	panicUtils "github.com/hdget/utils/panic"
	reflectUtils "github.com/hdget/utils/reflect"
	"github.com/pkg/errors"
)

type GinModule interface {
	Module
	RegisterHandlers(functions map[string]gin.HandlerFunc) error // 注册Handlers
	GetHandlers() map[string]gin.HandlerFunc                     // 获取handlers
}

type ginModuleImpl struct {
	Module
	self     any                        // 实际module实例
	handlers map[string]gin.HandlerFunc // 别名->handler
}

var (
	_ GinModule = (*ginModuleImpl)(nil)
)

func NewGinModule(moduleObject any, app string, functions map[string]gin.HandlerFunc) error {
	// 首先实例化module
	module, err := asGinModule(moduleObject, app)
	if err != nil {
		return err
	}

	// 然后注册handlers
	err = module.RegisterHandlers(functions)
	if err != nil {
		return err
	}

	// 最后注册module
	registerModule(module)

	return nil
}

// RegisterHandlers 参数handlers为gin.HandlerFunc
func (impl *ginModuleImpl) RegisterHandlers(functions map[string]gin.HandlerFunc) error {
	impl.handlers = make(map[string]gin.HandlerFunc)
	for alias, fn := range functions {
		impl.handlers[alias] = impl.newGinHandler(fn)
	}
	return nil
}

func (impl *ginModuleImpl) GetHandlers() map[string]gin.HandlerFunc {
	return impl.handlers
}

// asGinModule 将一个any类型的结构体转换成InvocationModule
//
// e,g:
//
//		type testModule struct {
//		  GinModule
//		}
//
//		 v := &testModule{}
//		 im, err := asGinModule(v, "test")
//	     if err != nil {
//	      ...
//	     }
func asGinModule(moduleObject any, app string) (GinModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &ginModuleImpl{
		Module: m,
		self:   moduleObject,
	}

	// 初始化module
	err = reflectUtils.StructSet(moduleObject, (*GinModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "devops module: %+v", m)
	}

	module, ok := moduleObject.(GinModule)
	if !ok {
		return nil, errors.New("invalid invocation module")
	}

	return module, nil
}

func (impl *ginModuleImpl) newGinHandler(fn gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				panicUtils.RecordErrorStack(impl.GetApp())
			}
		}()

		fn(ctx)
	}
}
