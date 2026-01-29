package loader

import (
	"bytes"
	"github.com/hdget/sdk/providers/config/viper/param"
	"github.com/spf13/viper"
)

type cliConfigLoader struct {
	localViper *viper.Viper
	param      *param.Cli
}

func NewCliConfigLoader(localViper *viper.Viper, param *param.Cli) Loader {
	return &cliConfigLoader{
		localViper: localViper,
		param:      param,
	}
}

// Load 从环境变量中读取配置信息
func (loader *cliConfigLoader) Load() error {
	if loader.param == nil {
		return nil
	}

	// 如果指定了配置内容，则合并
	if loader.param.Content != nil {
		_ = loader.localViper.MergeConfig(bytes.NewReader(loader.param.Content))
	}
	return nil
}
