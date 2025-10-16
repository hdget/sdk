package dapr

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/dapr/go-sdk/service/http"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/common/biz"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/common/types"
	"github.com/hdget/sdk/dapr/api"
	"github.com/hdget/sdk/dapr/module"
	"github.com/pkg/errors"
)

type hookPoint int

const (
	hookPointUnknown hookPoint = iota
	hookPointPreStart
	hookPointPreStop
	fileExposedHandlers = ".exposed_handlers.json"
)

type daprServerImpl struct {
	common.Service
	ctx    context.Context
	cancel context.CancelFunc
	debug  bool
	// 自定义参数
	app              string                             // 运行的app
	hooks            map[hookPoint][]types.HookFunction // 钩子函数
	registerFunction RegisterFunction                   // 向系统注册appServer的函数
	registerHandlers []*protobuf.DaprHandler            // 向系统注册的方法
	assets           embed.FS                           // 嵌入文件系统
	logger           types.LoggerProvider
	mq               types.MessageQueueProvider
}

func GetInvocationModules() []module.InvocationModule {
	return module.Get[module.InvocationModule](module.ModuleKindInvocation)
}

func NewGrpcServer(app, address string, options ...ServerOption) (types.AppServer, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("grpc server failed to listen on %s: %w", address, err)
	}

	// devops health check handler
	grpcServer := grpc.NewServiceWithListener(lis)

	ctx, cancel := context.WithCancel(context.Background())
	appServer := &daprServerImpl{
		Service: grpcServer,
		ctx:     ctx,
		cancel:  cancel,
		hooks:   make(map[hookPoint][]types.HookFunction),
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

func NewHttpServer(app, address string, options ...ServerOption) (types.AppServer, error) {
	httpServer := http.NewServiceWithMux(address, nil)

	ctx, cancel := context.WithCancel(context.Background())
	appServer := &daprServerImpl{
		Service: httpServer,
		ctx:     ctx,
		cancel:  cancel,
		hooks:   make(map[hookPoint][]types.HookFunction),
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
	for _, fn := range impl.hooks[hookPointPreStart] {
		if err := fn(); err != nil {
			return err
		}
	}

	impl.setupPreStopNotifier()

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
	impl.onPreStop()

	if len(forced) > 0 && forced[0] {
		return impl.Service.Stop()
	}
	return impl.GracefulStop()
}

func (impl *daprServerImpl) HookPreStart(hookFunctions ...types.HookFunction) types.AppServer {
	impl.hook(hookPointPreStart, hookFunctions...)
	return impl
}

func (impl *daprServerImpl) HookPreStop(hookFunctions ...types.HookFunction) types.AppServer {
	impl.hook(hookPointPreStop, hookFunctions...)
	return impl
}

// /////////////////////////////////////////////////////////////////////
// private functions
// /////////////////////////////////////////////////////////////////////
func (impl *daprServerImpl) hook(hookPoint hookPoint, hookFunctions ...types.HookFunction) {
	impl.hooks[hookPoint] = append(impl.hooks[hookPoint], hookFunctions...)
}

// Initialize 初始化server
func (impl *daprServerImpl) initialize() error {
	if err := impl.addHealthCheckHandler(); err != nil {
		return errors.Wrap(err, "adding health check handler")
	}

	if err := impl.addEventHandlers(); err != nil {
		return errors.Wrap(err, "adding event handlers")
	}

	if err := impl.subscribeDelayEvents(); err != nil {
		return errors.Wrap(err, "subscribe delay events")
	}

	if err := impl.addInvocationHandlers(); err != nil {
		return errors.Wrap(err, "adding invocation handlers")
	}

	return nil
}

// addEventHandlers 添加事件处理函数
func (impl *daprServerImpl) addEventHandlers() error {
	if impl.logger == nil {
		return errors.New("logger provider not found")
	}

	for _, m := range module.Get[module.EventModule](module.ModuleKindEvent) {
		for _, h := range m.GetHandlers() {
			e := api.NewEvent(m.GetPubSub(), h.GetTopic(), h.GetEventFunction(impl.logger))
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
	topic2delayEventHandler := make(map[string]module.DelayEventHandler)
	for _, m := range module.Get[module.DelayEventModule](module.ModuleKindDelayEvent) {
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

	healthModules := module.Get[module.HealthModule](module.ModuleKindHealth)
	if len(healthModules) == 0 {
		h = module.EmptyHealthCheckFunction
	} else {
		h = pie.First(healthModules).GetHandler()
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
	for _, m := range GetInvocationModules() {
		for _, h := range m.GetHandlers() {
			if impl.debug {
				impl.logger.Debug("add invocation handler", "method", h.GetInvokeName())
			}
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

	_, err := api.New(biz.NewContext()).Invoke("gateway", 1, "route", "update", &protobuf.UpdateRouteRequest{
		App:      app,
		Handlers: exposedHandlers,
	})
	if err != nil {
		return err
	}
	return nil
}

func (impl *daprServerImpl) setupPreStopNotifier() {
	// 监听中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGTERM, // 系统终止信号
	)

	go func() {
		<-sigCh
		impl.onPreStop()
	}()
}

func (impl *daprServerImpl) onPreStop() {
	for _, fn := range impl.hooks[hookPointPreStop] {
		if err := fn(); err != nil {
			if impl.logger != nil {
				impl.logger.Error("pre stop", "err", err)
			}
		}
	}

	impl.cancel()
}

// LoadStoredExposedHandlers 从embed.FS中加载ast解析后保存的DaprHandlers
func LoadStoredExposedHandlers(fs embed.FS) ([]*protobuf.DaprHandler, error) {
	// IMPORTANT: embedfs使用的是斜杠来获取文件路径,在windows平台下如果使用filepath来处理路径会导致问题
	data, err := fs.ReadFile(path.Join("json", fileExposedHandlers))
	if err != nil {
		return nil, err
	}

	var handlers []*protobuf.DaprHandler
	err = json.Unmarshal(data, &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}
