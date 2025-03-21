package sourcecode

import "github.com/hdget/sdk/dapr"

type Option func(*sourceCodeManagerImpl)

// WithServerCallSignature 定义服务调用函数签名， 方便定位哪个文件需要patch来添加模块所在包的导入路径
func WithServerCallSignature(importPath, functionChain string) Option {
	return func(m *sourceCodeManagerImpl) {
		m.serverCallSignature = &CallSignature{
			importPath:    importPath,
			functionChain: functionChain,
		}
	}
}

// WithHandlerMatchers 定义dapr module中handler匹配规则
func WithHandlerMatchers(matchers ...dapr.HandlerMatcher) Option {
	return func(m *sourceCodeManagerImpl) {
		m.handlerMatchers = matchers
	}
}

func WithSkipDirs(dirs ...string) Option {
	return func(m *sourceCodeManagerImpl) {
		m.skipDirs = append(m.skipDirs, dirs...)
	}
}
