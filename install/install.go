package install

import (
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

type TableInitFunction func(ctx biz.Context, fs embed.FS) error

type Installer interface {
	InstallDatabase(dbExecutor types.DbExecutor) (string, error)
	InstallTables(dbExecutor types.DbExecutor, force bool, tableNames ...string) error
}

type DbInstaller struct {
	createSQL    string
	initFunction TableInitFunction
}

type appInstallerImpl struct {
	project            string
	app                string
	assetFs            embed.FS
	tableInitFunctions map[string]TableInitFunction
}

const (
	psqlDropTable              = `DROP TABLE IF EXISTS "public"."%s";`
	psqlCreateDatabase         = `CREATE DATABASE %s WITH LC_COLLATE = 'C' LC_CTYPE = 'en_US.utf8' TABLESPACE = pg_default;`
	psqlClearPrepareStatements = `
DEALLOCATE ALL;       -- 清除当前会话的所有预处理语句
DISCARD PLANS;        -- PostgreSQL ≥13 替代方案`
)

func New(assetFs embed.FS, name string, tableInitFunctions map[string]TableInitFunction) (Installer, error) {
	project, exists := os.LookupEnv(constant.EnvKeyNamespace)
	if !exists || project == "" {
		return nil, fmt.Errorf("project name not found in %s", constant.EnvKeyNamespace)
	}

	return &appInstallerImpl{
		project:            project,
		app:                name,
		assetFs:            assetFs,
		tableInitFunctions: tableInitFunctions,
	}, nil
}

func (impl *appInstallerImpl) InstallDatabase(dbExecutor types.DbExecutor) (string, error) {
	dbName := fmt.Sprintf("%s_%s", impl.project, impl.app)

	sql := fmt.Sprintf(psqlCreateDatabase, dbName)

	_, err := dbExecutor.Exec(sql)
	if err != nil {
		return "", errors.Wrap(err, "create database")
	}

	return dbName, nil
}

func (impl *appInstallerImpl) InstallTables(dbExecutor types.DbExecutor, force bool, tableNames ...string) error {
	var sqlDir string
	dbKind := strings.Split(sdk.Db().GetCapability().Name, "-")[0]
	switch dbKind {
	case "postgresql":
		sqlDir = path.Join("sql", "postgresql")

	default:
		return fmt.Errorf("database type: %s not supported yet", dbKind)
	}

	// 清除所有预处理语句
	_, err := dbExecutor.Exec(psqlClearPrepareStatements)
	if err != nil {
		return err
	}

	tableName2sqlCreate, err := impl.findTableCreateSQL(sqlDir)
	if err != nil {
		return err
	}

	// 获取要处理的表
	installTables := tableNames
	if len(installTables) == 0 {
		installTables = pie.Keys(tableName2sqlCreate)
	}

	for _, tableName := range installTables {
		fmt.Printf("=== install %s table ===\n", tableName)
		if force {
			fmt.Printf(" * drop table: %s\n", tableName)
			sqlDrop := fmt.Sprintf(psqlDropTable, tableName)
			_, err = dbExecutor.Exec(sqlDrop)
			if err != nil {
				return err
			}
		}

		// create table
		if sqlCreate, exists := tableName2sqlCreate[tableName]; exists {
			fmt.Printf(" * create table: %s\n", tableName)
			_, err = dbExecutor.Exec(sqlCreate)
			if err != nil {
				return err
			}
		}

		// init table
		ctx := biz.NewContext()
		ctx.SetTx(dbExecutor)

		if tableInitFunction, exists := impl.tableInitFunctions[tableName]; exists {
			fmt.Printf(" * init table: %s\n", tableName)
			if err = tableInitFunction(ctx, impl.assetFs); err != nil {
				return err
			}
		}
	}
	return nil
}

func (impl *appInstallerImpl) findTableCreateSQL(dir string) (map[string]string, error) {
	entries, err := impl.assetFs.ReadDir(dir)
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

			sqlData, err := impl.assetFs.ReadFile(path.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}

			table2sql[tableName] = string(sqlData)
		}
	}

	return table2sql, nil
}
