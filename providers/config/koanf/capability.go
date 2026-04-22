package koanf

import (
	"github.com/hdget/sdk/common/provider"
	"go.uber.org/fx"
)

const (
	providerName = "config-koanf"
)

var Capability = provider.Capability{
	Category: provider.CategoryConfig,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
