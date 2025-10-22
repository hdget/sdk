package module

import (
	"github.com/hdget/common/namespace"
	reflectUtils "github.com/hdget/utils/reflect"
	"github.com/pkg/errors"
	"time"
)

type EventModule interface {
	Module
	RegisterHandlers(functions map[string]EventFunction) error // 注册Handlers
	GetHandlers() []eventHandler                               // 获取handlers
	GetPubSub() string
	GetAckTimeout() time.Duration
}

type eventModuleImpl struct {
	Module
	pubsub     string // 消息中间件名称定义在dapr配置中
	handlers   []eventHandler
	ackTimeout time.Duration
}

var (
	_                 EventModule = (*eventModuleImpl)(nil)
	defaultAckTimeout             = 29 * time.Minute // rabbitmq的默认超时时间为30分钟这里设置为29分钟保持
)

// NewEventModule 新建事件模块会执行下列操作:
func NewEventModule(moduleObject any, app, pubsub string, functions map[string]EventFunction, options ...EventModuleOption) error {
	// 首先实例化module
	module, err := asEventModule(moduleObject, app, namespace.Encapsulate(pubsub), options...)
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

// asEventModule 将一个any类型的结构体转换成EventModule
//
// e,g:
//
//		type v1_test struct {
//		  InvocationModule
//		}
//
//		 v := &v1_test{}
//		 im, err := asEventModule("app",v)
//	     if err != nil {
//	      ...
//	     }
//	     im.DiscoverHandlers()
func asEventModule(moduleObject any, app, pubsub string, options ...EventModuleOption) (EventModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &eventModuleImpl{
		Module:     m,
		pubsub:     pubsub,
		ackTimeout: defaultAckTimeout,
	}

	for _, option := range options {
		option(moduleInstance)
	}

	// 初始化module
	err = reflectUtils.StructSet(moduleObject, (*EventModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "devops module: %+v", m)
	}

	module, ok := moduleObject.(EventModule)
	if !ok {
		return nil, errors.New("invalid event module")
	}

	return module, nil
}

func (m *eventModuleImpl) GetKind() ModuleKind {
	return ModuleKindEvent
}

// RegisterHandlers 参数handlers为alias=>receiver.fnName, 保存为handler.id=>*invocationHandler
func (m *eventModuleImpl) RegisterHandlers(functions map[string]EventFunction) error {
	m.handlers = make([]eventHandler, 0)
	for topic, fn := range functions {
		m.handlers = append(m.handlers, m.newEventHandler(m, topic, fn))
	}
	return nil
}

func (m *eventModuleImpl) GetHandlers() []eventHandler {
	return m.handlers
}

func (m *eventModuleImpl) GetPubSub() string {
	return m.pubsub
}

func (m *eventModuleImpl) GetAckTimeout() time.Duration {
	return m.ackTimeout
}

func (m *eventModuleImpl) newEventHandler(module EventModule, topic string, fn EventFunction) eventHandler {
	return &eventHandlerImpl{
		module: module,
		topic:  topic,
		fn:     fn,
	}
}
