package sdk

import "github.com/hdget/common/types"

//type SdkConfig struct {
//	App            string
//	Debug          bool // debug mode
//	ConfigFilePath string
//	ConfigRootDirs []string
//}

type SdkOption func(*types.SdkConfig)

//
//var (
//	defaultConfigRootDirs = []string{
//		filepath.Join("setting", "app"),          // todo: old config root dir
//		filepath.Join("config", "app"),           // new config root dir
//		filepath.Join("common", "config", "app"), // match git directory
//	}
//	defaultEnvPrefix  = "HD"
//	defaultConfigType = "toml"
//)

func newSdkConfig(app string) *types.SdkConfig {
	return &types.SdkConfig{
		App:   app,
		Debug: false,
	}
}

func WithDebug() SdkOption {
	return func(c *types.SdkConfig) {
		c.Debug = true
	}
}

func WithConfigFile(configFilePath string) SdkOption {
	return func(c *types.SdkConfig) {
		c.ConfigFilePath = configFilePath
	}
}

func WithConfigRootDir(rootDirs ...string) SdkOption {
	return func(c *types.SdkConfig) {
		c.ConfigRootDirs = rootDirs
	}
}

//
//func (c *SdkConfig) getConfigProviderOption() *types.ConfigProviderOption {
//	configOption := &types.ConfigProviderOption{
//		App:             c.app,
//		ConfigEnvPrefix: defaultEnvPrefix,
//		ConfigRootDirs:  defaultConfigRootDirs, // 其他环境的BaseDir
//		ConfigType:      defaultConfigType,
//	}
//
//	if c.configFilePath != "" {
//		configOption.ConfigFile = c.configFilePath
//	}
//
//	if len(c.configRootDirs) > 0 {
//		configOption.ConfigRootDirs = c.configRootDirs
//	} else {
//		configOption.ConfigRootDirs = defaultConfigRootDirs
//	}
//
//	return configOption
//}
