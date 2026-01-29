package sdk

import (
	"github.com/hdget/sdk/providers/config/koanf"
)

type Option func(instance *SdkInstance)

func WithDebug() Option {
	return func(instance *SdkInstance) {
		instance.debug = true
	}
}

func WithConfigFile(configFile string) Option {
	return func(instance *SdkInstance) {
		instance.configOptions = append(instance.configOptions, koanf.WithConfigFile(configFile))
	}
}

func WithConfigContent(configContent []byte) Option {
	return func(instance *SdkInstance) {
		instance.configOptions = append(instance.configOptions, koanf.WithConfigContent(configContent))
	}
}
