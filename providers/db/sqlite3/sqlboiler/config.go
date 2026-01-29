package sqlite3_sqlboiler

import (
	"github.com/hdget/sdk/common/types"
	"github.com/pkg/errors"
)

type sqliteProviderConfig struct {
	DbPath string `mapstructure:"db"`
}

const (
	configSection = "sdk.sqlite"
)

var (
	errInvalidConfig = errors.New("invalid config")
	errEmptyConfig   = errors.New("empty config")
)

func newConfig(configProvider types.ConfigProvider) (*sqliteProviderConfig, error) {
	if configProvider == nil {
		return nil, errInvalidConfig
	}

	var c *sqliteProviderConfig
	err := configProvider.Unmarshal(&c, configSection)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errEmptyConfig
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate sqlite3 config")
	}

	return c, nil
}

func (c *sqliteProviderConfig) validate() error {
	if c == nil || c.DbPath == "" {
		return errInvalidConfig
	}
	return nil
}
