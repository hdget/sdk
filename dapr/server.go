package dapr

import (
	"context"
	"embed"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/dapr/go-sdk/service/http"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/common/intf"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/common/types"
	"github.com/pkg/errors"
	"net"
)

type daprServerImpl struct {
	common.Service
	ctx    context.Context
	cancel context.CancelFunc
	// 自定义参数
	app              string                                 // 运行的app
	hooks            map[intf.HookPoint][]intf.HookFunction // 钩子函数
	providers        []intf.Provider                        // sdk的providers
	registerFunction RegisterFunction                       // 向系统注册appServer的函数
	registerHandlers []*protobuf.DaprHandler                // 向系统注册的方法
	assets           embed.FS                               // 嵌入文件系统
	logger           intf.LoggerProvider
	mq               intf.MessageQueueProvider
}

var (
	_invocationModules = make([]InvocationModule, 0) // service invocation module
	_eventModules      = make([]EventModule, 0)      // topic event module
	_healthModules     = make([]HealthModule, 0)     // health module
	_delayEventModules = make([]DelayEventModule, 0) // delay event module
)

func GetInvocationModules() []InvocationModule {
	return _invocationModules
}

func NewGrpcServer(app, address string, options ...ServerOption) (intf.AppServer, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("grpc server failed to listen on %s: %w", address, err)
	}

	// install health check handler
	grpcServer := grpc.NewServiceWithListener(lis)

	ctx, cancel := context.WithCancel(context.Background())
	appServer := &daprServerImpl{
		Service: grpcServer,
		ctx:     ctx,
		cancel:  cancel,
		hooks:   make(map[intf.HookPoint][]intf.HookFunction),
		app:     app,
	}

	for _, apply := range options {
		apply(appServer)
	}

	if err = appServer.initialize(); err != nil {
		return nil, err
	}

	return appServer, nil
}

func NewHttpServer(app, address string, options ...ServerOption) (intf.AppServer, error) {
	httpServer := http.NewServiceWithMux(address, nil)

	ctx, cancel := context.WithCancel(context.Background())
	appServer := &daprServerImpl{
		Service: httpServer,
		ctx:     ctx,
		cancel:  cancel,
		hooks:   make(map[intf.HookPoint][]intf.HookFunction),
		app:     app,
	}

	for _, apply := range options {
		apply(appServer)
	}

	if err := appServer.initialize(); err != nil {
		return nil, err
	}

	return appServer, nil
}

func (impl *daprServerImpl) Start() error {
	for _, fn := range impl.hooks[intf.HookPointBeforeStart] {
		if err := fn(); err != nil {
			return err
		}
	}

	// 启动前调用AppRegister函数
	fnRegister := impl.registerFunction
	if fnRegister == nil {
		fnRegister = impl.defaultRegisterFunction
	}

	if err := fnRegister(impl.app, impl.registerHandlers); err != nil {
		return err
	}

	return impl.Service.Start()
}

func (impl *daprServerImpl) Stop(forced ...bool) error {
	for _, fn := range impl.hooks[intf.HookPointBeforeStop] {
		if err := fn(); err != nil {
			return err
		}
	}

	impl.cancel()
	if len(forced) > 0 && forced[0] {
		return impl.Service.Stop()
	}
	return impl.Service.GracefulStop()
}

func (impl *daprServerImpl) AddHook(hookPoint intf.HookPoint, hookFunctions ...intf.HookFunction) intf.AppServer {
	impl.hooks[hookPoint] = append(impl.hooks[hookPoint], hookFunctions...)
	return impl
}

///////////////////////////////////////////////////////////////////////
// private functions
///////////////////////////////////////////////////////////////////////

// Initialize 初始化server
func (impl *daprServerImpl) initialize() error {
	for _, provider := range impl.providers {
		switch provider.GetCapability().Category {
		case types.ProviderCategoryLogger:
			impl.logger = provider.(intf.LoggerProvider)
		case types.ProviderCategoryMq:
			impl.mq = provider.(intf.MessageQueueProvider)
		}
	}

	if err := impl.addHealthCheckHandler(); err != nil {
		return errors.Wrap(err, "adding health check handler")
	}

	if err := impl.addInvocationHandlers(); err != nil {
		return errors.Wrap(err, "adding invocation handlers")
	}

	if err := impl.addEventHandlers(); err != nil {
		return errors.Wrap(err, "adding event handlers")
	}

	if err := impl.subscribeDelayEvents(); err != nil {
		return errors.Wrap(err, "subscribe delay events")
	}

	return nil
}

// addEventHandlers 添加事件处理函数
func (impl *daprServerImpl) addEventHandlers() error {
	if impl.logger == nil {
		return errors.New("logger provider not found")
	}

	for _, m := range _eventModules {
		for _, h := range m.GetHandlers() {
			e := newEvent(m.GetPubSub(), h.GetTopic(), h.GetEventFunction(impl.logger))
			if err := impl.AddTopicEventHandler(e.Subscription, e.Handler); err != nil {
				return err
			}
		}
	}
	return nil
}

// subscribeDelayEvents 添加延迟事件处理函数
func (impl *daprServerImpl) subscribeDelayEvents() error {
	var app string
	topic2delayEventHandler := make(map[string]DelayEventHandler)
	for _, m := range _delayEventModules {
		if app == "" {
			app = m.GetApp()
		}

		for _, h := range m.GetHandlers() {
			topic2delayEventHandler[h.GetTopic()] = h
		}
	}

	if len(topic2delayEventHandler) == 0 {
		return nil
	}

	// delaySubscriber.SubscribeDelay需要指定subscribe的name
	if app == "" {
		return errors.New("app not found")
	}

	if impl.logger == nil {
		return errors.New("logger provider not found")
	}

	if impl.mq == nil {
		return errors.New("message queue provider not found")
	}

	delaySubscriber, err := impl.mq.NewSubscriber(app, &types.SubscriberOption{SubscribeDelayMessage: true})
	if err != nil {
		return errors.Wrapf(err, "new delay event subscriber, name: %s", app)
	}

	for _, h := range topic2delayEventHandler {
		msgChan, err := delaySubscriber.Subscribe(impl.ctx, h.GetTopic())
		if err != nil {
			return errors.Wrapf(err, "subscribe topic, topic: %s", h.GetTopic())
		}

		impl.logger.Debug("subscribe delay event", "topic", h.GetTopic())
		go h.Handle(impl.ctx, impl.logger, msgChan)
	}
	return nil
}

// addHealthCheckHandler 添加健康检测Handler
func (impl *daprServerImpl) addHealthCheckHandler() error {
	var h common.HealthCheckHandler
	if len(_healthModules) == 0 {
		h = EmptyHealthCheckFunction
	} else {
		h = pie.First(_healthModules).GetHandler()
	}

	// 注册health check handler
	return impl.AddHealthCheckHandler("", h)
}

// addHealthCheckHandler 添加服务调用Handler
func (impl *daprServerImpl) addInvocationHandlers() error {
	if impl.logger == nil {
		return errors.New("logger provider not found")
	}

	// 注册各种类型的handlers
	for _, m := range _invocationModules {
		for _, h := range m.GetHandlers() {
			if err := impl.AddServiceInvocationHandler(h.GetInvokeName(), h.GetInvokeFunction(impl.logger)); err != nil {
				return err
			}
		}
	}
	return nil
}

// defaultRegisterFunction 缺省的将AppServer注册到系统的函数
func (impl *daprServerImpl) defaultRegisterFunction(app string, handlers []*protobuf.DaprHandler) error {
	exposedHandlers := handlers

	// 如果没有传入handlers, 尝试从资源目录的.exposed_handlers.json中加载handlers
	if len(exposedHandlers) == 0 {
		var err error
		exposedHandlers, err = LoadStoredExposedHandlers(impl.assets)
		if err != nil && impl.logger != nil {
			impl.logger.Debug("load exposed handlers", "err", err)
		}
	}

	if len(exposedHandlers) == 0 {
		return nil
	}

	_, err := Api().Invoke("gateway", 1, "route", "update", &protobuf.UpdateRouteRequest{
		App:      app,
		Handlers: handlers,
	})
	if err != nil {
		return err
	}
	return nil
}

func registerModule(module any) {
	switch m := module.(type) {
	case InvocationModule:
		_invocationModules = append(_invocationModules, m)
	case EventModule:
		_eventModules = append(_eventModules, m)
	case DelayEventModule:
		_delayEventModules = append(_delayEventModules, m)
	case HealthModule:
		_healthModules = append(_healthModules, m)
	}
}
