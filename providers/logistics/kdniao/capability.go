package logistics_kdniao

import (
	"github.com/hdget/sdk/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "logistics-kdniao"
)

var Capability = types.Capability{
	Category: types.ProviderCategoryLogistics,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
