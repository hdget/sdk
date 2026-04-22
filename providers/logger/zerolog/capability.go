package zerolog

import (
	"github.com/hdget/sdk/common/provider"
	"go.uber.org/fx"
)

const (
	providerName = "logger-zerolog"
)

var Capability = provider.Capability{
	Category: provider.CategoryLogger,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
