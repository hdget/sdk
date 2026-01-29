package loader

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
)

const (
	// 最小化的配置,保证日志工作正常
	tplMinimalConfigContent = `
[sdk]
  [sdk.logger]
	   level = "debug"
	   filename = "%s.log"
	   [sdk.logger.rotate]
		   max_age = 7`
)

type minimalConfigLoader struct {
	app        string
	localViper *viper.Viper
}

func NewMinimalConfigLoader(app string, localViper *viper.Viper) Loader {
	return &minimalConfigLoader{
		app:        app,
		localViper: localViper,
	}
}

func (loader *minimalConfigLoader) Load() error {
	minimalConfig := fmt.Sprintf(tplMinimalConfigContent, loader.app)
	return loader.localViper.MergeConfig(bytes.NewReader([]byte(minimalConfig)))
}
