package sdk

import (
	viper "github.com/hdget/provider-config-viper"
)

type Option func(instance *SdkInstance)

func WithDebug() Option {
	return func(instance *SdkInstance) {
		instance.debug = true
	}
}

func WithConfigFile(configFile string) Option {
	return func(instance *SdkInstance) {
		instance.configOptions = append(instance.configOptions, viper.WithConfigFile(configFile))
	}
}

func WithConfigContent(configContent []byte) Option {
	return func(instance *SdkInstance) {
		instance.configOptions = append(instance.configOptions, viper.WithConfigContent(configContent))
	}
}

func WithDefaultRemote() Option {
	return func(instance *SdkInstance) {
		instance.configOptions = append(instance.configOptions, viper.WithDefaultRemote())
	}
}
