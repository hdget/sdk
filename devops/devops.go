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
	"github.com/pkg/errors"
)

type TableOperator interface {
	GetName() string
	Init(ctx biz.Context, fs embed.FS) error
	Export(ctx biz.Context, assetPath string) error
}

type DevOps interface {
	InstallDatabase(dbname ...string) (string, error)
	InstallTables(ctx biz.Context, store embed.FS, force bool, tableNames ...string) error
	ExportTables(ctx biz.Context, storePath string, tableNames ...string) error
}

type devOpsImpl struct {
	app               string
	db                types.DbProvider
	dbKind            string
	tableOperators    []TableOperator
	needDangerConfirm bool
}

const (
	dbKindPostgreSQL = "postgresql"
	dbKindMySQL      = "mysql"
	dbKindSqlite3    = "sqlite3"
)

var (
	supportedDbKinds = []string{dbKindPostgreSQL, dbKindMySQL, dbKindSqlite3}
)

func New(app string, db types.DbProvider, options ...Option) (DevOps, error) {
	if db == nil {
		return nil, fmt.Errorf("db provider not provided")
	}

	dbKind := strings.Split(db.GetCapability().Name, "-")[0]
	if !pie.Contains(supportedDbKinds, dbKind) {
		return nil, fmt.Errorf("db kind %s not supported", dbKind)
	}

	impl := &devOpsImpl{
		app:               app,
		db:                db,
		dbKind:            dbKind,
		tableOperators:    nil,
		needDangerConfirm: true,
	}

	for _, option := range options {
		option(impl)
	}

	return impl, nil
}

func (impl *devOpsImpl) InstallDatabase(dbName ...string) (string, error) {
	project, exists := os.LookupEnv(constant.EnvKeyNamespace)
	if !exists || project == "" {
		return "", fmt.Errorf("project name not found in %s", constant.EnvKeyNamespace)
	}

	databaseName := fmt.Sprintf("%s_%s", project, impl.app)
	// if specified dbname, use db
	if len(dbName) > 0 && dbName[0] != "" {
		databaseName = dbName[0]
	}
	fmt.Printf("=== install database: %s ===\n", databaseName)

	switch impl.dbKind {
	case dbKindPostgreSQL:
		sql := fmt.Sprintf(psqlCreateDatabase, databaseName)
		_, err := impl.db.My().Exec(sql)
		if err != nil {
			return "", errors.Wrap(err, "create database")
		}
	case dbKindSqlite3:
		fmt.Println("PLEASE use below command to create database:")
		fmt.Println()
		fmt.Printf("sqlite3 %s\n", databaseName)
		fmt.Println()
		return "", errors.New("automatically create database not supported")
	}

	return databaseName, nil
}

func (impl *devOpsImpl) InstallTables(ctx biz.Context, store embed.FS, force bool, tableNames ...string) error {
	tx, ok := ctx.Transactor().GetTx().(types.DbExecutor)
	if !ok {
		return fmt.Errorf("db transactor not found in context")
	}

	// 清除所有预处理语句
	switch impl.dbKind {
	case dbKindPostgreSQL:
		_, err := tx.Exec(psqlBeforeInstallTables)
		if err != nil {
			return err
		}
	}

	// 获取SQL文件
	tableName2sqlCreate, err := impl.findTableCreateSQL(store, path.Join("sql", impl.dbKind))
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
			var sqlDrop string
			switch impl.dbKind {
			case dbKindPostgreSQL:
				sqlDrop = fmt.Sprintf(psqlDropTable, tableName)
			case dbKindSqlite3:
				sqlDrop = fmt.Sprintf(sqlite3DropTable, tableName)
			}

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
