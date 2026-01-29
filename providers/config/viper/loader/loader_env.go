package loader

import (
	"github.com/hdget/providers/config/viper/param"
	"github.com/spf13/viper"
)

type envLoader struct {
	localViper *viper.Viper
	param      *param.Env
}

func NewEnvConfigLoader(localViper *viper.Viper, param *param.Env) Loader {
	return &envLoader{
		localViper: localViper,
		param:      param,
	}
}

// Load 从环境变量中读取配置信息
func (loader *envLoader) Load() error {
	if loader.param == nil {
		return nil
	}

	loader.localViper.SetEnvPrefix(loader.param.Prefix)
	loader.localViper.AutomaticEnv()
	return nil
}
