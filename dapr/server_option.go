package dapr

import (
	"embed"
	"github.com/hdget/common/intf"
	"github.com/hdget/common/protobuf"
)

// RegisterFunction app向gateway注册的函数
type RegisterFunction func(string, []*protobuf.DaprHandler) error

type ServerOption func(impl *daprServerImpl)

// WithProviders 提供的providers
func WithProviders(providers ...intf.Provider) ServerOption {
	return func(impl *daprServerImpl) {
		impl.providers = append(impl.providers, providers...)
	}
}

func WithRegisterFunction(fn RegisterFunction) ServerOption {
	return func(impl *daprServerImpl) {
		impl.registerFunction = fn
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
