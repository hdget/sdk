package viper

import (
	"github.com/hdget/sdk/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "config-viper"
)

var Capability = types.Capability{
	Category: types.ProviderCategoryConfig,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
