package logistics_kdniao

import (
	"github.com/hdget/sdk/common/types"
	"github.com/pkg/errors"
)

type kdniaoProviderConfig struct {
	Default *kdniaoClientConfig   `mapstructure:"default"`
	Items   []*kdniaoClientConfig `mapstructure:"items"`
}

type kdniaoClientConfig struct {
	Name        string `mapstructure:"name"`
	EBusinessID string `mapstructure:"ebusiness_id"`
	AppKey      string `mapstructure:"app_key"`
}

const (
	configSection = "sdk.logistics.kdniao"
)

var (
	errInvalidConfig = errors.New("invalid kdniao config")
	errEmptyConfig   = errors.New("empty kdniao config")
)

func newConfig(configProvider types.ConfigProvider) (*kdniaoProviderConfig, error) {
	if configProvider == nil {
		return nil, errInvalidConfig
	}

	var c *kdniaoProviderConfig
	err := configProvider.Unmarshal(&c, configSection)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errEmptyConfig
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate kdniao provider config")
	}

	return c, nil
}

func (c *kdniaoProviderConfig) validate() error {
	if c.Default != nil {
		err := c.validateClientConfig(c.Default)
		if err != nil {
			return err
		}
	}

	for _, item := range c.Items {
		err := c.validateExtraClientConfig(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *kdniaoProviderConfig) validateClientConfig(conf *kdniaoClientConfig) error {
	if conf.EBusinessID == "" || conf.AppKey == "" {
		return errEmptyConfig
	}
	return nil
}

func (c *kdniaoProviderConfig) validateExtraClientConfig(conf *kdniaoClientConfig) error {
	if conf.Name == "" || conf.EBusinessID == "" || conf.AppKey == "" {
		return errInvalidConfig
	}
	return nil
}
