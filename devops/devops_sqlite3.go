package devops

import (
	"embed"
	"fmt"
	"path"

	"github.com/elliotchance/pie/v2"
	"github.com/hdget/common/biz"
	"github.com/hdget/common/types"
	"github.com/pkg/errors"
)

type sqlite3DevOpsImpl struct {
	*devOpsImpl
}

const (
	sqlite3DropTable = `DROP TABLE IF EXISTS %s;`
	sqlite3StoreDir  = "sqlite3"
)

func SQLite3(app string, options ...Option) DevOps {
	return &sqlite3DevOpsImpl{
		devOpsImpl: newDevOps(app, options...),
	}
}

func (impl *sqlite3DevOpsImpl) InstallDatabase(dbClient types.DbClient, specifiedDbName ...string) (string, error) {
	dbName, err := impl.getDbName(specifiedDbName...)
	if err != nil {
		return "", errors.Wrap(err, "get db name")
	}

	fmt.Printf("=== install database: %s ===\n", dbName)

	// do nothing, just close db client
	_ = dbClient.Close()

	return "", nil
}

func (impl *sqlite3DevOpsImpl) InstallTables(ctx biz.Context, store embed.FS, force bool, tableNames ...string) error {
	tx, ok := ctx.Transactor().GetTx().(types.DbExecutor)
	if !ok {
		return fmt.Errorf("db transactor not found in context")
	}

	// 获取SQL文件
	tableName2sqlCreate, err := impl.findTableCreateSQL(store, path.Join("sql", sqlite3StoreDir))
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

			sqlDrop := fmt.Sprintf(sqlite3DropTable, tableName)
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
