package dapr

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/common/protobuf"
)

// AppRegisterFunction app向gateway注册的函数
type AppRegisterFunction func([]*protobuf.DaprHandler) error

type ServerOption func(impl *daprServerImpl)

// WithProviders 提供的providers
func WithProviders(providers ...intf.Provider) ServerOption {
	return func(impl *daprServerImpl) {
		impl.providers = append(impl.providers, providers...)
	}
}

func WithGatewayRegisterFunction(fn AppRegisterFunction) ServerOption {
	return func(impl *daprServerImpl) {
		impl.fnAppRegister = fn
	}
}

func WithInvocationHandlers(handlers []*protobuf.DaprHandler) ServerOption {
	return func(impl *daprServerImpl) {
		impl.invocationHandlers = handlers
	}
}
