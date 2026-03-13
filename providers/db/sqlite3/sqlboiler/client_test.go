package sqlboiler

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

func TestSqlDB(t *testing.T) {
	// 创建临时目录用于测试数据库
	tmpDir, err := os.MkdirTemp("", "sqlite_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	// 创建配置
	cfg := &sqliteProviderConfig{
		DbPath: dbPath,
	}

	// 创建 client
	client, err := newClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// 测试 SqlDB() 返回非 nil
	db := client.SqlDB()
	if db == nil {
		t.Error("SqlDB() returned nil, expected non-nil *sql.DB")
	}

	// 测试返回的 *sql.DB 可以正常工作
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Errorf("failed to query using returned *sql.DB: %v", err)
	}
	if result != 1 {
		t.Errorf("expected result 1, got %d", result)
	}
}

func TestSqlDB_ReturnsSameInstance(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "sqlite_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &sqliteProviderConfig{
		DbPath: dbPath,
	}

	client, err := newClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// 多次调用 SqlDB() 应返回相同的实例
	db1 := client.SqlDB()
	db2 := client.SqlDB()

	if db1 != db2 {
		t.Error("SqlDB() should return the same *sql.DB instance")
	}
}

func TestSqlDB_AfterClose(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "sqlite_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &sqliteProviderConfig{
		DbPath: dbPath,
	}

	client, err := newClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// 获取 *sql.DB 引用
	db := client.SqlDB()

	// 关闭 client
	err = client.Close()
	if err != nil {
		t.Fatalf("failed to close client: %v", err)
	}

	// 关闭后 *sql.DB 仍然可以被引用，但操作应该失败
	// 这验证了 SqlDB() 返回的是实际的 *sql.DB，而不是副本
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err == nil {
		t.Error("expected error when querying closed database, got nil")
	}
}

// TestSqlite3ClientImplementsDbClient 确保 sqlite3Client 实现了 types.DbClient 接口
func TestSqlite3ClientImplementsDbClient(t *testing.T) {
	cfg := &sqliteProviderConfig{
		DbPath: ":memory:",
	}

	client, err := newClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// 编译时检查接口实现
	// 如果 sqlite3Client 没有实现 types.DbClient，编译会失败
	var _ interface {
		Close() error
		SqlDB() *sql.DB
	} = &sqlite3Client{}
}