package redigo

import (
	"github.com/hdget/sdk/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "redis-redigo"
)

var Capability = types.Capability{
	Category: types.ProviderCategoryRedis,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
