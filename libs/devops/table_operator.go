package devops

import (
	"context"
	"embed"
)

type tableOperatorImpl struct {
	name string
}

func NewTableOperator(name string) TableOperator {
	return &tableOperatorImpl{
		name: name,
	}
}

func (impl *tableOperatorImpl) Init(ctx context.Context, fs embed.FS) error {
	return nil
}

func (impl *tableOperatorImpl) Export(ctx context.Context, assetPath string) error {
	return nil
}

func (impl *tableOperatorImpl) GetName() string {
	return impl.name
}