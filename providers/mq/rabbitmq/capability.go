package rabbitmq

import (
	"github.com/hdget/sdk/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "mq-rabbitmq"
)

var Capability = types.Capability{
	Category: types.ProviderCategoryMq,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
