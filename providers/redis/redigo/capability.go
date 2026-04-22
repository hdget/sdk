package redigo

import (
	"github.com/hdget/sdk/common/provider"
	"go.uber.org/fx"
)

const (
	providerName = "redis-redigo"
)

var Capability = provider.Capability{
	Category: provider.CategoryRedis,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
