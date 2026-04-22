package rabbitmq

import (
	"github.com/hdget/sdk/common/provider"
	"go.uber.org/fx"
)

const (
	providerName = "mq-rabbitmq"
)

var Capability = provider.Capability{
	Category: provider.CategoryMq,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
