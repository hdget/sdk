package sqlc

import (
	"github.com/hdget/sdk/common/provider"
	"go.uber.org/fx"
)

const (
	providerName = "mysql-sqlc"
)

var Capability = provider.Capability{
	Category: provider.CategoryDb,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}