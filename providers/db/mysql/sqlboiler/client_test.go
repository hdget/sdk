package sqlboiler

import (
	"database/sql"
	"testing"
)

// TestMysqlClientDb_MethodExists 测试 Db 方法存在且返回正确类型
func TestMysqlClientDb_MethodExists(t *testing.T) {
	// 这个测试验证 mysqlClient 结构体有 Db 方法
	// 编译时检查：如果方法不存在或签名错误，编译会失败
	client := &mysqlClient{DB: nil}

	// 调用 Db 应该返回 nil（因为我们设置 DB 为 nil）
	db := client.Db()
	if db != nil {
		t.Error("expected nil from Db() when DB is nil")
	}
}

// TestMysqlClientImplementsDbClient 确保 mysqlClient 实现了必要的接口
func TestMysqlClientImplementsDbClient(t *testing.T) {
	// 编译时检查接口实现
	var _ interface {
		Close() error
		Db() *sql.DB
	} = &mysqlClient{}
}