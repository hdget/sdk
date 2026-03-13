package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/hdget/sdk/common/types"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type psqlClient struct {
	*sql.DB
}

const (
	// DSN (Data Type NickName): username:password@address/dbname?param=value
	dsnTemplate = "postgres://%s:%s@%s:%d/%s?TimeZone=Asia/Shanghai"
)

func newClient(c *psqlConfig) (types.DbClient, error) {
	// 构造连接参数
	dsn := fmt.Sprintf(dsnTemplate, c.User, url.QueryEscape(c.Password), c.Host, c.Port, c.Database)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxOpenConn)                     // 最大活跃连接数（对应PgBouncer的max_client_conn）
	db.SetMaxIdleConns(int(0.2 * float32(c.MaxOpenConn))) // 最大空闲连接数（建议值=0.2*MaxOpenConns)
	if !c.UsePgBouncer {
		db.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second) // 应小于PostgreSQL的idle_in_transaction_session_timeout, 默认为60m
	}

	return &psqlClient{DB: db}, nil
}

func (m *psqlClient) Close() error {
	return m.DB.Close()
}

func (m *psqlClient) SqlDB() *sql.DB {
	return m.DB
}

// RunInTransaction 在事务中执行函数，支持嵌套事务（通过 SAVEPOINT 实现）
func (m *psqlClient) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
