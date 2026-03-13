package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hdget/sdk/common/types"
	_ "modernc.org/sqlite"
)

type sqlite3Client struct {
	*sql.DB
}

const (
	dsnTemplate = "file:%s?_loc=Local"
)

func newClient(c *sqliteProviderConfig, args ...string) (types.DbClient, error) {
	var absDbFile string
	if len(args) > 0 {
		absDbFile = args[0]
	} else {
		if !filepath.IsAbs(c.DbPath) {
			workDir, _ := os.Getwd()
			absDbFile = filepath.Join(workDir, c.DbPath)
		} else {
			absDbFile = c.DbPath
		}
	}

	// 构造连接参数
	dsn := fmt.Sprintf(dsnTemplate, absDbFile)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	var userVersion int
	err = db.QueryRow("PRAGMA user_version").Scan(&userVersion)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("fail connect db: %s", absDbFile)
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	// https://making.pusher.com/production-ready-connection-pooling-in-go
	// Avoid issue:
	// packets.go:123: closing bad idle connection: EOF
	// connection.go:173: driver: bad connection
	db.SetConnMaxLifetime(3 * time.Minute)

	return &sqlite3Client{DB: db}, nil
}

func (m *sqlite3Client) Close() error {
	return m.DB.Close()
}

func (m *sqlite3Client) SqlDB() *sql.DB {
	return m.DB
}

// RunInTransaction 在事务中执行函数，支持嵌套事务（通过 SAVEPOINT 实现）
func (m *sqlite3Client) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// 检查是否已在事务中
	if tx, ok := ctx.Value(types.TxCtxKey{}).(*sql.Tx); ok {
		// 已在事务中，创建 SAVEPOINT 实现嵌套事务
		spName := fmt.Sprintf("sp_%d", time.Now().UnixNano())
		_, _ = tx.ExecContext(ctx, fmt.Sprintf("SAVEPOINT %s", spName))

		err := fn(ctx)
		if err != nil {
			_, _ = tx.ExecContext(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", spName))
			return err
		}
		_, _ = tx.ExecContext(ctx, fmt.Sprintf("RELEASE SAVEPOINT %s", spName))
		return nil
	}

	// 开始新事务
	tx, err := m.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 将事务放入 context，支持嵌套事务检测
	txCtx := context.WithValue(ctx, types.TxCtxKey{}, tx)
	err = fn(txCtx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
