package sdk

import viper "github.com/hdget/provider-config-viper"

type Option func(instance *SdkInstance)

func WithDebug() Option {
	return func(instance *SdkInstance) {
		instance.debug = true
	}
}

func WithConfigOptions(configOptions ...viper.Option) Option {
	return func(instance *SdkInstance) {
		for _, option := range configOptions {
			option(instance.configParam)
		}
	}
}
