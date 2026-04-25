package dapr

import (
	"embed"

	"github.com/hdget/sdk/common/protobuf"
	"github.com/hdget/sdk/common/provider"
)

// RegisterFunction app向gateway注册的函数
type RegisterFunction func(string, []*protobuf.DaprHandler) error

type ServerOption func(impl *daprServerImpl)

// WithProviders 提供的providers
func WithProviders(providers ...provider.Provider) ServerOption {
	return func(impl *daprServerImpl) {
		// try initialize provider
		for _, p := range providers {
			switch p.GetCapability().Category {
			case provider.CategoryLogger:
				impl.logger = p.(provider.Logger)
			case provider.CategoryMq:
				impl.mq = p.(provider.MessageQueue)
			}
		}
	}
}

func WithRegisterFunction(fn RegisterFunction) ServerOption {
	return func(impl *daprServerImpl) {
		impl.registerFunction = fn
	}
}

func WithSkipRegister() ServerOption {
	return func(impl *daprServerImpl) {
		impl.registerFunction = func(string, []*protobuf.DaprHandler) error {
			return nil
		}
	}
}

func WithRegisterHandlers(handlers []*protobuf.DaprHandler) ServerOption {
	return func(impl *daprServerImpl) {
		impl.registerHandlers = handlers
	}
}

func WithAssets(fs embed.FS) ServerOption {
	return func(impl *daprServerImpl) {
		impl.assets = fs
	}
}

func WithDebug(debug bool) ServerOption {
	return func(impl *daprServerImpl) {
		impl.debug = debug
	}
}
