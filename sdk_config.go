package sdk

import "github.com/hdget/common/types"

type Option func(*types.SdkConfig)

func newSdkConfig(app string) *types.SdkConfig {
	return &types.SdkConfig{
		App:   app,
		Debug: false,
	}
}

func WithDebug() Option {
	return func(c *types.SdkConfig) {
		c.Debug = true
	}
}

func WithConfigFile(configFilePath string) Option {
	return func(c *types.SdkConfig) {
		c.ConfigFilePath = configFilePath
	}
}

func WithConfigRootDir(rootDirs ...string) Option {
	return func(c *types.SdkConfig) {
		c.ConfigRootDirs = rootDirs
	}
}
