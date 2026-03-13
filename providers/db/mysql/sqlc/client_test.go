package sqlc

import (
	"database/sql"
	"testing"
)

// TestMysqlClientSqlDB_MethodExists 测试 SqlDB 方法存在且返回正确类型
func TestMysqlClientSqlDB_MethodExists(t *testing.T) {
	// 这个测试验证 mysqlClient 结构体有 SqlDB 方法
	client := &mysqlClient{DB: nil}

	// 调用 SqlDB 应该返回 nil（因为我们设置 DB 为 nil）
	db := client.SqlDB()
	if db != nil {
		t.Error("expected nil from SqlDB() when DB is nil")
	}
}

// TestMysqlClientImplementsDbClient 确保 mysqlClient 实现了必要的接口
func TestMysqlClientImplementsDbClient(t *testing.T) {
	// 编译时检查接口实现
	var _ interface {
		Close() error
		SqlDB() *sql.DB
	} = &mysqlClient{}
}