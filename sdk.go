package sdk

import (
	"context"
	"github.com/hdget/common/intf"
	"github.com/hdget/common/types"
	"github.com/hdget/provider-config-viper"
	"github.com/hdget/provider-logger-zerolog"
	"github.com/hdget/utils/logger"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"sync"
)

type SdkInstance struct {
	config *types.SdkConfig

	configProvider intf.ConfigProvider
	loggerProvider intf.LoggerProvider
	dbProvider     intf.DbProvider
	redisProvider  intf.RedisProvider
	mqProvider     intf.MessageQueueProvider
	ossProvider    intf.OssProvider

	configVar any // 配置数据
}

var (
	_instance *SdkInstance
	once      sync.Once

	errInvalidCapability = errors.New("invalid capability")
)

func New(app string, options ...Option) *SdkInstance {
	once.Do(
		func() {
			sdkConfig := newSdkConfig(app)
			for _, apply := range options {
				apply(sdkConfig)
			}
			_instance = &SdkInstance{
				config: sdkConfig,
			}
		},
	)
	return _instance
}

func HasInitialized() bool {
	return _instance != nil
}

func GetInstance() *SdkInstance {
	return _instance
}

// UseConfig 加载配置信息到给定的配置变量中。
// 该方法使用配置提供者（configProvider）将配置数据解析到配置变量中。
// 参数:
//
//	configVar - 一个配置变量的指针，用于接收解析后的配置数据。
func (i *SdkInstance) UseConfig(configVar any) *SdkInstance {
	i.configVar = configVar
	return i
}

// Initialize initializes the SDK instance with given capabilities.
// This function configures the SDK instance using dependency injection with fx.Options,
// based on the provided capabilities, such as database, logging, and configuration providers.
func (i *SdkInstance) Initialize(capabilities ...types.Capability) error {
	// Prepare fxOptions for DI configuration
	fxOptions := []fx.Option{
		fx.Provide(func() *types.SdkConfig { return i.config }),
		viper.Capability.Module, // Initialize configProvider
		fx.Populate(&_instance.configProvider),
		zerolog.Capability.Module, // Initialize loggerProvider
		fx.Populate(&_instance.loggerProvider),
	}

	// Iterate through the provided capabilities and configure the corresponding providers
	for _, c := range capabilities {
		switch c.Category {
		case types.ProviderCategoryDb:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.dbProvider))
		case types.ProviderCategoryRedis:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.redisProvider))
		case types.ProviderCategoryMq:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.mqProvider))
		case types.ProviderCategoryOss:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.ossProvider))
		default:
			return errors.Wrapf(errInvalidCapability, "capability: %s", c.Name)
		}
	}

	// Disable fx internal logger in production mode
	if !i.config.Debug {
		fxOptions = append(fxOptions, fx.NopLogger)
	}

	// Start the DI container and initialize all configured providers
	err := fx.New(fxOptions...).Start(context.Background())
	if err != nil {
		return errors.Unwrap(err)
	}

	// try load config
	i.loadConfig()

	return nil
}

func (i *SdkInstance) loadConfig() {
	// 如果没有赋值，则直接返回
	if i.configVar == nil {
		return
	}
	// 检查配置提供者是否已初始化。
	if i.configProvider == nil {
		// 如果未初始化，则记录致命错误并终止程序。
		logger.Fatal("config provider not initialized")
	}

	// 将配置数据解析到配置变量中。
	err := i.configProvider.Unmarshal(i.configVar)
	if err != nil {
		// 如果解析失败，则记录致命错误并终止程序。
		logger.Fatal("unmarshal to config variable", "err", err)
	}
}
