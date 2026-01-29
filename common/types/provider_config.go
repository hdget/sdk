package types

type ConfigProvider interface {
	Provider
	Unmarshal(configVar any, key ...string) error // 读取配置到变量configVar
	Get(key string) any                           // 获取配置项的值
}
