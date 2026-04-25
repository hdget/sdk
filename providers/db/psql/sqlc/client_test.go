package sqlc

import (
	"database/sql"
	"testing"
)

// TestPsqlClientDb_MethodExists 测试 Db 方法存在且返回正确类型
func TestPsqlClientDb_MethodExists(t *testing.T) {
	// 这个测试验证 psqlClient 结构体有 Db 方法
	client := &psqlClient{DB: nil}

	// 调用 Db 应该返回 nil（因为我们设置 DB 为 nil）
	db := client.Db()
	if db != nil {
		t.Error("expected nil from Db() when DB is nil")
	}
}

// TestPsqlClientImplementsDbClient 确保 psqlClient 实现了必要的接口
func TestPsqlClientImplementsDbClient(t *testing.T) {
	// 编译时检查接口实现
	var _ interface {
		Close() error
		Db() *sql.DB
	} = &psqlClient{}
}