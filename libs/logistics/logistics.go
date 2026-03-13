package logistics

import "sync"

// Factory 创建物流API实例的工厂函数类型
type Factory func(cfg *Config) (LogisticsApi, error)

var (
	factoryRegistry = make(map[string]Factory)
	registryMutex   sync.RWMutex
)

// RegisterFactory 注册物流API工厂函数
// 通常在实现包的 init() 中调用
func RegisterFactory(name string, factory Factory) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	factoryRegistry[name] = factory
}

// New 根据配置创建物流API实例
// 支持 kdniao 和 kd100 两种供应商
func New(cfg *Config) (LogisticsApi, error) {
	if cfg == nil {
		return nil, ErrEmptyConfig
	}

	registryMutex.RLock()
	factory, ok := factoryRegistry[cfg.Name]
	registryMutex.RUnlock()

	if !ok {
		return nil, ErrUnknownVendor
	}

	return factory(cfg)
}
