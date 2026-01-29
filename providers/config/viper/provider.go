package viper

import (
	"os"

	"github.com/hdget/sdk/common/constant"
	"github.com/hdget/sdk/common/types"
	"github.com/hdget/sdk/providers/config/viper/loader"
	"github.com/hdget/sdk/providers/config/viper/param"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// viperConfigProvider 命令行配置
type viperConfigProvider struct {
	app   string
	env   string
	local *viper.Viper
	param *param.Param
}

// New 初始化config provider
func New(app string, options ...Option) (types.ConfigProvider, error) {
	env, exists := os.LookupEnv(constant.EnvKeyRunEnvironment)
	if !exists {
		return nil, errors.New("env not found")
	}

	provider := &viperConfigProvider{
		app:   app,
		env:   env,
		local: viper.New(),
		param: param.GetDefaultParam(),
	}

	for _, option := range options {
		option(provider.param)
	}

	err := provider.Load()
	if err != nil {
		return nil, errors.Wrap(err, "load config")
	}

	return provider, nil
}

// Load 从各个配置源获取配置数据, 并加载到configVar中，同名变量优先级高的覆盖低的
// - remote: kvstore配置（低）
// - configFile: 文件配置(中)
// - env: 环境变量配置(高)
// - input: 命令行参数配置(最高)
func (p *viperConfigProvider) Load() error {
	// 如果环境变量为空，则加载最小基本配置
	if p.env == "" {
		return loader.NewMinimalConfigLoader(p.app, p.local).Load()
	}

	// 尝试从环境变量中获取配置信息
	if err := loader.NewFileConfigLoader(p.local, p.param.File, p.app, p.env).Load(); err != nil {
		return errors.Wrap(err, "load config from file")
	}

	// 尝试从环境变量中获取配置信息
	if err := loader.NewEnvConfigLoader(p.local, p.param.Env).Load(); err != nil {
		return errors.Wrap(err, "load config from env")
	}

	// 尝试从环境变量中获取配置信息
	if err := loader.NewCliConfigLoader(p.local, p.param.Cli).Load(); err != nil {
		return errors.Wrap(err, "load config from cli")
	}

	// 尝试从远程配置中获取配置信息
	if err := loader.NewRemoteConfigLoader(p.local, p.param.Remote, p.env).Load(); err != nil {
		return errors.Wrap(err, "load config from remote")
	}

	return nil
}

// Unmarshal 解析配置
func (p *viperConfigProvider) Unmarshal(configVar any, args ...string) error {
	if len(args) > 0 {
		return p.local.UnmarshalKey(args[0], configVar)
	}
	return p.local.Unmarshal(configVar)
}

func (p *viperConfigProvider) Get(key string) any {
	return p.local.Get(key)
}

func (p *viperConfigProvider) GetCapability() types.Capability {
	return Capability
}
