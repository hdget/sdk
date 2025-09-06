package dapr

import (
	"embed"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/common/types"
)

// RegisterFunction app向gateway注册的函数
type RegisterFunction func(string, []*protobuf.DaprHandler) error

type ServerOption func(impl *daprServerImpl)

// WithProviders 提供的providers
func WithProviders(providers ...types.Provider) ServerOption {
	return func(impl *daprServerImpl) {
		// try initialize provider
		for _, provider := range providers {
			switch provider.GetCapability().Category {
			case types.ProviderCategoryLogger:
				impl.logger = provider.(types.LoggerProvider)
			case types.ProviderCategoryMq:
				impl.mq = provider.(types.MessageQueueProvider)
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
