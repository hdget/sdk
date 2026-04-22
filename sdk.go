package sdk

import (
	"context"
	"sync"

	"github.com/hdget/sdk/common/provider"
	"github.com/hdget/sdk/providers/config/koanf"
	"github.com/hdget/sdk/providers/logger/zerolog"
	"github.com/hdget/utils/logger"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type Instance struct {
	configProvider provider.Config
	loggerProvider provider.Logger
	dbProvider     provider.Database
	redisProvider  provider.Redis
	mqProvider     provider.MessageQueue
	app            string
	debug          bool
	configVar      any            // 配置变量
	configOptions  []koanf.Option // 配置选项
}

var (
	_instance                *Instance
	once                     sync.Once
	errUnsupportedCapability = errors.New("unsupported capability")
)

func New(app string, options ...Option) *Instance {
	once.Do(
		func() {
			_instance = &Instance{
				app:           app,
				configOptions: make([]koanf.Option, 0),
			}

			for _, apply := range options {
				apply(_instance)
			}

		},
	)
	return _instance
}

func HasInitialized() bool {
	return _instance != nil
}

func GetInstance() *Instance {
	return _instance
}

// UseConfig 加载配置信息到给定的配置变量中。
// 该方法使用配置提供者（configProvider）将配置数据解析到配置变量中。
// 参数:
//
//	configVar - 一个配置变量的指针，用于接收解析后的配置数据。
func (i *Instance) UseConfig(configVar any) *Instance {
	i.configVar = configVar
	return i
}

// Initialize initializes the SDK instance with given capabilities.
// This function configures the SDK instance using dependency injection with fx.Options,
// based on the provided capabilities, such as database, logging, and configuration providers.
func (i *Instance) Initialize(capabilities ...provider.Capability) error {
	// Prepare fxOptions for DI configuration
	fxOptions := []fx.Option{
		// Initialize configProvider
		fx.Provide(func() (provider.Config, error) {
			return koanf.New(i.app, i.configOptions...)
		}),
		fx.Populate(&_instance.configProvider),
		// Initialize loggerProvider
		zerolog.Capability.Module,
		fx.Populate(&_instance.loggerProvider),
	}

	// Iterate through the provided capabilities and configure the corresponding providers
	for _, c := range capabilities {
		switch c.Category {
		case provider.CategoryDb:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.dbProvider))
		case provider.CategoryRedis:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.redisProvider))
		case provider.CategoryMq:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.mqProvider))
		default:
			return errors.Wrapf(errUnsupportedCapability, "capability: %s", c.Name)
		}
	}

	// Disable fx internal logger in production mode
	if !i.debug {
		fxOptions = append(fxOptions, fx.NopLogger)
	}

	// Start the DI container and initialize all configured providers
	err := fx.New(fxOptions...).Start(context.Background())
	if err != nil {
		return errors.Unwrap(err)
	}

	// try load config to config var
	i.unmarshalConfig()

	return nil
}

func (i *Instance) unmarshalConfig() {
	var fatal, outputError func(msg string, keyvals ...interface{})

	if i.loggerProvider != nil {
		fatal = i.loggerProvider.Fatal
		outputError = i.loggerProvider.Error
	} else {
		fatal = logger.Fatal
		outputError = logger.Error
	}

	// 检查配置提供者是否已初始化。
	if i.configProvider == nil {
		// 如果未初始化，则记录致命错误并终止程序。
		fatal("config provider not initialized")
	}

	// 如果没有赋值，则直接返回
	if i.configVar != nil {
		// 将配置数据解析到局部配置变量中
		err := i.configProvider.Unmarshal(i.configVar)
		if err != nil {
			// 如果解析失败，则记录致命错误并终止程序。
			outputError("unmarshal to config variable", "err", err)
		}
	}
}
