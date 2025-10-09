package devops

import (
	"embed"

	"github.com/hdget/common/biz"
)

type tableOperatorImpl struct {
	name string
}

func NewTableOperator(name string) TableOperator {
	return &tableOperatorImpl{
		name: name,
	}
}

func (impl *tableOperatorImpl) Init(ctx biz.Context, fs embed.FS) error {
	return nil
}

func (impl *tableOperatorImpl) Export(ctx biz.Context, fs embed.FS) error {
	return nil
}

func (impl *tableOperatorImpl) GetName() string {
	return impl.name
}
