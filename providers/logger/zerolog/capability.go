package zerolog

import (
	"github.com/hdget/sdk/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "logger-zerolog"
)

var Capability = types.Capability{
	Category: types.ProviderCategoryLogger,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
