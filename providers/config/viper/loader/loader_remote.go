package loader

import (
	"fmt"
	"github.com/hdget/sdk/common/constant"
	"github.com/hdget/sdk/providers/config/viper/param"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"strings"
	"time"
)

type remoteConfigLoader struct {
	viper *viper.Viper
	param *param.Remote
	env   string
}

func NewRemoteConfigLoader(viper *viper.Viper, param *param.Remote, env string) Loader {
	return &remoteConfigLoader{
		viper: viper,
		param: param,
		env:   env,
	}
}

// Load 从环境变量中读取配置信息
func (loader *remoteConfigLoader) Load() error {
	if loader.param == nil || loader.param.Provider == "" {
		return nil
	}

	// 如果endpoints为空，尝试已读取的配置中获取
	if len(loader.param.Endpoints) == 0 {
		configEndpoints := loader.viper.GetStringSlice(fmt.Sprintf(constant.ConfigKeyRemoteEndpoints, loader.param.Provider))
		if len(configEndpoints) != 0 {
			loader.param.Endpoints = configEndpoints
		} else {
			loader.param.Endpoints = param.DefaultRemoteEndpoints
		}
	}

	// 如果watchPath为空，设置缺省值
	if loader.param.WatchPath == "" {
		loader.param.WatchPath = fmt.Sprintf(constant.ConfigPathRemote, loader.env)
	}

	var err error
	if loader.param.Secret != "" {
		err = loader.viper.AddSecureRemoteProvider(
			loader.param.Provider,
			strings.Join(loader.param.Endpoints, ";"),
			loader.param.WatchPath,
			loader.param.Secret,
		)
	} else {
		err = loader.viper.AddRemoteProvider(
			loader.param.Provider,
			strings.Join(loader.param.Endpoints, ";"),
			loader.param.WatchPath,
		)
	}
	if err != nil {
		return errors.Wrap(err, "add remote Provider")
	}

	loader.viper.SetConfigType(loader.param.RemoteConfigType)

	// 尝试读取，不报错
	_ = loader.viper.ReadRemoteConfig()

	// 自动读取到kvstore
	if err = loader.viper.WatchRemoteConfigOnChannel(); err != nil {
		return err
	}

	// open a goroutine to unmarshal remote config
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(loader.param.WatchInterval))

			if loader.param.WatchCallback != nil {
				loader.param.WatchCallback()
			}
		}
	}()

	return nil
}
