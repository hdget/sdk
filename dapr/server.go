package dapr

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/dapr/go-sdk/service/http"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/common/intf"
	"github.com/hdget/common/types"
	"github.com/pkg/errors"
	"net"
)

type daprServerImpl struct {
	common.Service
	ctx    context.Context
	cancel context.CancelFunc
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

//func NewServer(assetFs embed.FS, options ...ServerOption) (intf.AppServer, error) {
//	// 解析go:embed路径
//	_, callFile, _, _ := runtime.Caller(1)
//	embedAbsPath, embedRelPath, err := newAstEmbedFinder(callFile).Parse()
//	if err != nil {
//		return nil, err
//	}
//
//	srv := &daprServerImpl{
//		assetManager:      asset.New(assetFs, embedAbsPath, embedRelPath),
//		actions:           make([]Action, 0),
//		serverImportPath:  defaultAppServerImportPath,
//		serverRunFuncName: defaultAppServerRunFunction,
//	}
//
//	for _, apply := range options {
//		apply(srv)
//	}
//	return srv, nil
//}

func NewGrpcServer(address string, providers ...intf.Provider) (intf.AppServer, error) {
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
	}

	if err = appServer.initialize(providers...); err != nil {
		return nil, err
	}

	return appServer, nil
}

func NewHttpServer(address string, providers ...intf.Provider) (intf.AppServer, error) {
	httpServer := http.NewServiceWithMux(address, nil)

	ctx, cancel := context.WithCancel(context.Background())
	appServer := &daprServerImpl{
		Service: httpServer,
		ctx:     ctx,
		cancel:  cancel,
	}

	if err := appServer.initialize(providers...); err != nil {
		return nil, err
	}

	return appServer, nil
}

func (impl *daprServerImpl) Start() error {
	return impl.Service.Start()
}

func (impl *daprServerImpl) Stop(forced ...bool) error {
	impl.cancel()
	if len(forced) > 0 && forced[0] {
		return impl.Service.Stop()
	}
	return impl.Service.GracefulStop()
}

///////////////////////////////////////////////////////////////////////
// private functions
///////////////////////////////////////////////////////////////////////

// Initialize 初始化server
func (impl *daprServerImpl) initialize(providers ...intf.Provider) error {
	var (
		loggerProvider intf.LoggerProvider
		mqProvider     intf.MessageQueueProvider
	)
	for _, provider := range providers {
		switch provider.GetCapability().Category {
		case types.ProviderCategoryLogger:
			loggerProvider = provider.(intf.LoggerProvider)
		case types.ProviderCategoryMq:
			mqProvider = provider.(intf.MessageQueueProvider)
		}
	}

	if err := impl.addHealthCheckHandler(); err != nil {
		return errors.Wrap(err, "adding health check handler")
	}

	if err := impl.addInvocationHandlers(loggerProvider); err != nil {
		return errors.Wrap(err, "adding invocation handlers")
	}

	if err := impl.addEventHandlers(loggerProvider); err != nil {
		return errors.Wrap(err, "adding event handlers")
	}

	if err := impl.subscribeDelayEvents(loggerProvider, mqProvider); err != nil {
		return errors.Wrap(err, "subscribe delay events")
	}

	return nil
}

// addEventHandlers 添加事件处理函数
func (impl *daprServerImpl) addEventHandlers(logger intf.LoggerProvider) error {
	if logger == nil {
		return errors.New("logger provider not found")
	}

	for _, m := range _eventModules {
		for _, h := range m.GetHandlers() {
			e := NewEvent(m.GetPubSub(), h.GetTopic(), h.GetEventFunction(logger))
			if err := impl.AddTopicEventHandler(e.Subscription, e.Handler); err != nil {
				return err
			}
		}
	}
	return nil
}

// subscribeDelayEvents 添加延迟事件处理函数
func (impl *daprServerImpl) subscribeDelayEvents(logger intf.LoggerProvider, mq intf.MessageQueueProvider) error {
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

	if logger == nil {
		return errors.New("logger provider not found")
	}

	if mq == nil {
		return errors.New("message queue provider not found")
	}

	delaySubscriber, err := mq.NewSubscriber(app, &types.SubscriberOption{SubscribeDelayMessage: true})
	if err != nil {
		return errors.Wrapf(err, "new delaySubscriber, name: %s", app)
	}

	for _, h := range topic2delayEventHandler {
		msgChan, err := delaySubscriber.Subscribe(impl.ctx, h.GetTopic())
		if err != nil {
			return errors.Wrapf(err, "subscribe topic, topic: %s", h.GetTopic())
		}

		logger.Debug("subscribe delay event", "topic", h.GetTopic())
		go h.Handle(impl.ctx, logger, msgChan)
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
func (impl *daprServerImpl) addInvocationHandlers(logger intf.LoggerProvider) error {
	if logger == nil {
		return errors.New("logger provider not found")
	}

	// 注册各种类型的handlers
	for _, m := range _invocationModules {
		for _, h := range m.GetHandlers() {
			if err := impl.AddServiceInvocationHandler(h.GetInvokeName(), h.GetInvokeFunction(logger)); err != nil {
				return err
			}
		}
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
