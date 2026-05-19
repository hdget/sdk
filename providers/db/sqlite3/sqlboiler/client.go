package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hdget/sdk/common/provider"
	_ "modernc.org/sqlite"
)

type sqlite3Client struct {
	*sql.DB
	logger provider.Logger
}

const (
	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	dsnTemplate = "file:%s?_loc=Local"
)

func newClient(c *sqliteProviderConfig, logger provider.Logger, args ...string) (provider.DbClient, error) {
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

	return &sqlite3Client{DB: db, logger: logger}, nil
}

func (m *sqlite3Client) Close() error {
	return m.DB.Close()
}

func (m *sqlite3Client) Db() *sql.DB {
	return m.DB
}

// RunInTransaction 在事务中执行函数，支持嵌套事务（通过 SAVEPOINT 实现）
func (m *sqlite3Client) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// 检查是否已在事务中
	if tx, ok := ctx.Value(provider.TxCtxKey{}).(*sql.Tx); ok {
		// 已在事务中，创建 SAVEPOINT 实现嵌套事务
		spName := fmt.Sprintf("sp_%d", time.Now().UnixNano())

		// 创建 SAVEPOINT
		_, err := tx.ExecContext(ctx, fmt.Sprintf("SAVEPOINT %s", spName))
		if err != nil {
			// SAVEPOINT 创建失败，说明事务可能已 aborted
			// 记录详细错误信息，帮助定位问题
			m.logger.Error("create savepoint failed",
				"savepoint", spName,
				"error", err,
				"hint", "transaction may be aborted by previous SQL error")
			return fmt.Errorf("create savepoint %s failed: %w", spName, err)
		}

		err = fn(ctx)
		if err != nil {
			// 回滚到 SAVEPOINT
			_, rbErr := tx.ExecContext(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", spName))
			if rbErr != nil {
				m.logger.Error("rollback to savepoint failed",
					"savepoint", spName,
					"error", rbErr,
					"original_error", err)
				// 返回原始错误，但记录回滚失败信息
				return fmt.Errorf("rollback to savepoint %s failed: %w (original: %v)", spName, rbErr, err)
			}
			return err
		}

		// 释放 SAVEPOINT
		_, relErr := tx.ExecContext(ctx, fmt.Sprintf("RELEASE SAVEPOINT %s", spName))
		if relErr != nil {
			m.logger.Warn("release savepoint failed",
				"savepoint", spName,
				"error", relErr,
				"hint", "savepoint will be released at transaction commit")
		}
		return nil
	}

	// 开始新事务
	tx, err := m.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 将事务放入 context，支持嵌套事务检测
	txCtx := context.WithValue(ctx, provider.TxCtxKey{}, tx)
	err = fn(txCtx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
