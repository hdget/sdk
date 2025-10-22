package devops

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/hdget/common/biz"
	"github.com/hdget/common/constant"
	"github.com/hdget/common/types"
	"github.com/hdget/sdk"
	"github.com/pkg/errors"
)

type devOpsImpl struct {
	project           string
	app               string
	tableOperators    []TableOperator
	needDangerConfirm bool
}

const (
	psqlDropTable              = `DROP TABLE IF EXISTS "public"."%s";`
	psqlCreateDatabase         = `CREATE DATABASE %s WITH LC_COLLATE = 'C' LC_CTYPE = 'en_US.utf8' TABLESPACE = pg_default;`
	psqlClearPrepareStatements = `
DEALLOCATE ALL;       -- 清除当前会话的所有预处理语句
DISCARD PLANS;        -- PostgreSQL ≥13 替代方案`
)

func New(name string, options ...Option) (Operator, error) {
	project, exists := os.LookupEnv(constant.EnvKeyNamespace)
	if !exists || project == "" {
		return nil, fmt.Errorf("project name not found in %s", constant.EnvKeyNamespace)
	}

	impl := &devOpsImpl{
		project:           project,
		app:               name,
		needDangerConfirm: true,
	}

	for _, option := range options {
		option(impl)
	}

	return impl, nil
}

func (impl *devOpsImpl) InstallDatabase(dbExecutor types.DbExecutor) (string, error) {
	dbName := fmt.Sprintf("%s_%s", impl.project, impl.app)
	fmt.Printf("=== install database: %s ===\n", dbName)

	sql := fmt.Sprintf(psqlCreateDatabase, dbName)

	_, err := dbExecutor.Exec(sql)
	if err != nil {
		return "", errors.Wrap(err, "create database")
	}

	return dbName, nil
}

func (impl *devOpsImpl) InstallTables(ctx biz.Context, store embed.FS, force bool, tableNames ...string) error {
	var sqlDir string
	dbKind := strings.Split(sdk.Db().GetCapability().Name, "-")[0]
	switch dbKind {
	case "postgresql":
		sqlDir = path.Join("sql", "postgresql")

	default:
		return fmt.Errorf("database type: %s not supported yet", dbKind)
	}

	tx, ok := ctx.Transactor().GetTx().(types.DbExecutor)
	if !ok {
		return fmt.Errorf("db transactor not found in context")
	}

	// 清除所有预处理语句
	_, err := tx.Exec(psqlClearPrepareStatements)
	if err != nil {
		return err
	}

	tableName2sqlCreate, err := impl.findTableCreateSQL(store, sqlDir)
	if err != nil {
		return err
	}

	// 获取要处理的表
	installTables := tableNames
	if len(installTables) == 0 {
		installTables = pie.Keys(tableName2sqlCreate)
	}

	for _, tableName := range installTables {
		fmt.Printf("=== install table: %s ===\n", tableName)
		if force {
			if impl.needDangerConfirm {
				prompt := fmt.Sprintf("WARNING: You are about to drop the table '%s'.\nThis action will permanently erase all data in the table and is IRREVERSIBLE!", tableName)
				confirmed, err := impl.confirm(prompt, "ok")
				if err != nil {
					return err
				}

				if !confirmed {
					continue
				}
			}

			fmt.Printf(" * drop table: %s\n", tableName)
			sqlDrop := fmt.Sprintf(psqlDropTable, tableName)
			_, err = tx.Exec(sqlDrop)
			if err != nil {
				return err
			}
		}

		// create table
		if sqlCreate, exists := tableName2sqlCreate[tableName]; exists {
			fmt.Printf(" * create table: %s\n", tableName)
			_, err = tx.Exec(sqlCreate)
			if err != nil {
				return err
			}
		}

		// init table
		foundIndex := pie.FindFirstUsing(impl.tableOperators, func(v TableOperator) bool {
			return v.GetName() == tableName
		})

		if foundIndex >= 0 {
			fmt.Printf(" * init table: %s\n", tableName)
			if err = impl.tableOperators[foundIndex].Init(ctx, store); err != nil {
				return err
			}
		}

	}
	return nil
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
