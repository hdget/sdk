package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/sdk/common/types"
)

type mysqlClient struct {
	*sql.DB
}

const (
	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	dsnTemplate = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"
)

func newClient(c *mysqlConfig) (types.DbClient, error) {
	// 构造连接参数
	dsn := fmt.Sprintf(dsnTemplate, c.User, c.Password, c.Host, c.Port, c.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	// https://making.pusher.com/production-ready-connection-pooling-in-go
	// Avoid issue:
	// packets.go:123: closing bad idle connection: EOF
	// connection.go:173: driver: bad connection
	db.SetConnMaxLifetime(3 * time.Minute)

	return &mysqlClient{DB: db}, nil
}

func (m *mysqlClient) Close() error {
	return m.DB.Close()
}

func (m *mysqlClient) SqlDB() *sql.DB {
	return m.DB
}

// RunInTransaction 在事务中执行函数，支持嵌套事务（通过 SAVEPOINT 实现）
func (m *mysqlClient) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
