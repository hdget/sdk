package devops

import (
	"embed"
	"fmt"

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
	fmt.Println(" * skipped")
	return nil
}

func (impl *tableOperatorImpl) Export(ctx biz.Context, assetPath string) error {
	fmt.Println(" * skipped")
	return nil
}

func (impl *tableOperatorImpl) GetName() string {
	return impl.name
}
