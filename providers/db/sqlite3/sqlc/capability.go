package sqlc

import (
	"github.com/hdget/sdk/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "sqlite3-sqlc"
)

var Capability = types.Capability{
	Name:     providerName,
	Category: types.ProviderCategoryDb,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}