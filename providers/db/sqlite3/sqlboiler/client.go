package sqlboiler

import (
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
	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
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

func (m sqlite3Client) Close() error {
	return m.DB.Close()
}

func (m sqlite3Client) Get(dest interface{}, query string, args ...interface{}) error {
	return nil
}

func (m sqlite3Client) Select(dest interface{}, query string, args ...interface{}) error {
	return nil
}

func (m sqlite3Client) Rebind(query string) string {
	return ""
}
