package module

import (
	"github.com/cenkalti/backoff/v4"
	reflectUtils "github.com/hdget/utils/reflect"
	"github.com/pkg/errors"
	"time"
)

type DelayEventModule interface {
	Module
	RegisterHandlers(functions map[string]DelayEventFunction) error // 注册Handlers
	GetHandlers() []DelayEventHandler                               // 获取handlers
	GetAckTimeout() time.Duration
	GetBackOffPolicy() backoff.BackOff
}

type delayEventModuleImpl struct {
	Module
	handlers      []DelayEventHandler
	ackTimeout    time.Duration
	backoffPolicy backoff.BackOff
}

var (
	_ DelayEventModule = (*delayEventModuleImpl)(nil)
)

// NewDelayEventModule 初始化延迟消息模块
func NewDelayEventModule(moduleObject any, app string, functions map[string]DelayEventFunction, options ...DelayEventModuleOption) error {
	// 首先实例化module
	module, err := asDelayEventModule(moduleObject, app, options...)
	if err != nil {
		return err
	}

	// 然后注册handlers
	err = module.RegisterHandlers(functions)
	if err != nil {
		return err
	}

	// 最后注册module
	register(module)

	return nil
}

// asDelayEventModule 将一个any类型的结构体转换成DelayEventModule
//
// e,g:
//
//		type v1_test struct {
//		  DelayEventModule
//		}
//
//		 v := &v1_test{}
//		 im, err := asDelayEventModule("app",v)
//	     if err != nil {
//	      ...
//	     }
//	     im.DiscoverHandlers()
func asDelayEventModule(moduleObject any, app string, options ...DelayEventModuleOption) (DelayEventModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &delayEventModuleImpl{
		Module:        m,
		ackTimeout:    defaultAckTimeout,
		backoffPolicy: getDefaultBackOffPolicy(),
	}

	for _, option := range options {
		option(moduleInstance)
	}

	// 初始化module
	err = reflectUtils.StructSet(moduleObject, (*DelayEventModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "devops module: %+v", m)
	}

	module, ok := moduleObject.(DelayEventModule)
	if !ok {
		return nil, errors.New("invalid delay event module")
	}

	return module, nil
}

func (impl *delayEventModuleImpl) GetKind() ModuleKind {
	return ModuleKindDelayEvent
}

// RegisterHandlers 参数handlers为alias=>receiver.fnName, 保存为handler.id=>*invocationHandler
func (impl *delayEventModuleImpl) RegisterHandlers(functions map[string]DelayEventFunction) error {
	impl.handlers = make([]DelayEventHandler, 0)
	for topic, fn := range functions {
		impl.handlers = append(impl.handlers, impl.newDelayEventHandler(impl, topic, fn))
	}
	return nil
}

func (impl *delayEventModuleImpl) GetHandlers() []DelayEventHandler {
	return impl.handlers
}

func (impl *delayEventModuleImpl) GetAckTimeout() time.Duration {
	return impl.ackTimeout
}

func (impl *delayEventModuleImpl) GetBackOffPolicy() backoff.BackOff {
	return impl.backoffPolicy
}

func (impl *delayEventModuleImpl) newDelayEventHandler(module DelayEventModule, topic string, fn DelayEventFunction) DelayEventHandler {
	return &delayEventHandlerImpl{
		module: module,
		topic:  topic,
		fn:     fn,
	}
}

// NewExponentialBackOff creates an instance of ExponentialBackOff using default values.
func getDefaultBackOffPolicy() backoff.BackOff {
	// 最开始等待3秒
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 3 * time.Second

	// 最多尝试3次
	nb := backoff.WithMaxRetries(b, 3)
	backoff.WithMaxRetries(b, 3) // 最多重试3次
	nb.Reset()
	return nb
}
