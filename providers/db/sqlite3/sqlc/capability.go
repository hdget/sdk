package sqlc

import (
	"github.com/hdget/sdk/common/provider"
	"go.uber.org/fx"
)

const (
	providerName = "sqlite3-sqlc"
)

var Capability = provider.Capability{
	Name:     providerName,
	Category: provider.CategoryDb,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}