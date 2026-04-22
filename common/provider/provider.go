package provider

import (
	"go.uber.org/fx"
)

// Provider 底层库能力提供者接口
type Provider interface {
	GetCapability() Capability // 获取能力
}

// Capability 能力提供者
type Capability struct {
	Category Category
	Name     string
	Module   fx.Option
}

type Category int

const (
	CategoryUnknown Category = iota
	CategoryConfig
	CategoryLogger
	CategoryDb
	CategoryRedis
	CategoryMq
)
