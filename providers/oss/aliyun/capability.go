package oss_aliyun

import (
	"github.com/hdget/sdk/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "oss-aliyun"
)

var Capability = types.Capability{
	Category: types.ProviderCategoryOss,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
