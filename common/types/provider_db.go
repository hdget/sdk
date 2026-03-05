package types

import (
	"context"
	"database/sql"
)

// TxCtxKey 用于在 context 中传递事务状态
type TxCtxKey struct{}

type DbProvider interface {
	Provider
	My() DbClient
	Master() DbClient
	Slave(i int) DbClient
	By(name string) DbClient
	// Read 返回用于读操作的数据库客户端（自动从 slave 中轮询选择）
	Read() DbClient
	// Write 返回用于写操作的数据库客户端（返回 master 或 default）
	Write() DbClient
}

type DbExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// DbContextExecutor can perform SQL queries with context
type DbContextExecutor interface {
	DbExecutor

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	PrepareContext(context.Context, string) (*sql.Stmt, error)
}

// DbClient 数据库客户端接口
type DbClient interface {
	DbContextExecutor

	Close() error
	// RunInTransaction 在事务中执行函数，支持嵌套事务（通过 SAVEPOINT 实现）
	// fn 的参数 ctx 包含事务信息，用于嵌套事务检测
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
