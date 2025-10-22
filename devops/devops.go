package devops

import (
	"embed"
	"github.com/hdget/common/biz"
	"github.com/hdget/common/types"
)

type TableOperator interface {
	GetName() string
	Init(ctx biz.Context, fs embed.FS) error
	Export(ctx biz.Context, assetPath string) error
}

type Operator interface {
	InstallDatabase(dbExecutor types.DbExecutor) (string, error)
	InstallTables(ctx biz.Context, store embed.FS, force bool, tableNames ...string) error
	ExportTables(ctx biz.Context, storePath string, tableNames ...string) error
}
