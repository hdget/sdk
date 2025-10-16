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
	InstallDatabase(executor types.DbExecutor) (string, error)
	InstallTables(executor types.DbExecutor, store embed.FS, force bool, tableNames ...string) error
	ExportTables(executor types.DbExecutor, storePath string, tableNames ...string) error
}
