package devops

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/hdget/sdk/common/biz"
	"github.com/hdget/sdk/common/constant"
	"github.com/hdget/sdk/common/types"
)

type TableOperator interface {
	GetName() string
	Init(ctx biz.Context, fs embed.FS) error
	Export(ctx biz.Context, assetPath string) error
}

type DevOps interface {
	InstallDatabase(dbClient types.DbClient, dbname ...string) (string, error)
	InstallTables(ctx biz.Context, store embed.FS, force bool, tableNames ...string) error
	ExportTables(ctx biz.Context, storePath string, tableNames ...string) error
}

type devOpsImpl struct {
	app               string
	tableOperators    []TableOperator
	needDangerConfirm bool
}

func newDevOps(app string, options ...Option) *devOpsImpl {
	impl := &devOpsImpl{
		app:               app,
		tableOperators:    nil,
		needDangerConfirm: true,
	}

	for _, option := range options {
		option(impl)
	}

	return impl
}

func (impl *devOpsImpl) getDbName(dbName ...string) (string, error) {
	if len(dbName) > 0 && dbName[0] != "" {
		return dbName[0], nil
	}

	project, exists := os.LookupEnv(constant.EnvKeyNamespace)
	if !exists || project == "" {
		return "", fmt.Errorf("project name not found in %s", constant.EnvKeyNamespace)
	}
	return fmt.Sprintf("%s_%s", project, impl.app), nil
}

func (impl *devOpsImpl) ExportTables(ctx biz.Context, storePath string, tableNames ...string) error {
	// 获取要处理的表
	exportTables := tableNames
	if len(exportTables) == 0 {
		exportTables = pie.Map(impl.tableOperators, func(v TableOperator) string {
			return v.GetName()
		})
	}

	for _, tableName := range exportTables {
		foundIndex := pie.FindFirstUsing(impl.tableOperators, func(v TableOperator) bool {
			return v.GetName() == tableName
		})

		if foundIndex >= 0 {
			fmt.Printf("=== export table: %s ===\n", tableName)
			if err := impl.tableOperators[foundIndex].Export(ctx, storePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (impl *devOpsImpl) findTableCreateSQL(fs embed.FS, dir string) (map[string]string, error) {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	table2sql := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasPrefix(entry.Name(), "table_") {
			tableName := strings.TrimSuffix(strings.TrimPrefix(entry.Name(), "table_"), ".sql")

			sqlData, err := fs.ReadFile(path.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}

			table2sql[tableName] = string(sqlData)
		}
	}

	return table2sql, nil
}

func (impl *devOpsImpl) confirm(prompt string, confirmAnswer string) (bool, error) {
	fmt.Printf("%s\n\nPlease type '%s' and press Enter to confirm deletion, or type anything else to cancel:", prompt, confirmAnswer)

	// 2. 读取用户输入
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("read input: %v", err)
	}

	// 3. 清理输入并检查确认信息
	input = strings.TrimSpace(input) // 去除输入字符串两端的空白字符（如回车符）
	if input != confirmAnswer {
		fmt.Println("operation cancelled")
		return false, nil
	}

	return true, nil
}
