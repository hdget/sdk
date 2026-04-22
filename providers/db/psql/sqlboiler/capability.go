package sqlboiler

import (
	"github.com/hdget/sdk/common/provider"
	"go.uber.org/fx"
)

const (
	providerName = "postgresql-sqlboiler"
)

var Capability = provider.Capability{
	Category: provider.CategoryDb,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(newProvider),
	),
}
