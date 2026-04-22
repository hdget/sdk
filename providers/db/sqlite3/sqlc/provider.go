package sqlc

import (
	"github.com/hdget/sdk/common/provider"
	"github.com/pkg/errors"
)

type sqlite3Provider struct {
	client provider.DbClient
}

func New(configProvider provider.Config, logger provider.Logger) (provider.Database, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	client, err := newClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "new sqlite3 client")
	}

	logger.Debug("init sqlite3 provider", "db", config.DbPath)
	return &sqlite3Provider{client: client}, nil
}

// NewClient 从指定的文件创建创建数据库连接
func NewClient(dbFile string) (provider.DbClient, error) {
	client, err := newClient(nil, dbFile)
	if err != nil {
		return nil, errors.Wrapf(err, "connect sqlite3: %s", dbFile)
	}
	return client, nil
}

func (p *sqlite3Provider) GetCapability() provider.Capability {
	return Capability
}

func (p *sqlite3Provider) Default() provider.DbClient {
	return p.client
}

func (p *sqlite3Provider) Master() provider.DbClient {
	return p.client
}

func (p *sqlite3Provider) Slave(i int) provider.DbClient {
	return p.client
}

func (p *sqlite3Provider) Named(name string) provider.DbClient {
	return p.client
}

// Read 返回用于读操作的数据库客户端（SQLite3 无读写分离，返回同一客户端）
func (p *sqlite3Provider) Read() provider.DbClient {
	return p.client
}

// Write 返回用于写操作的数据库客户端（SQLite3 无读写分离，返回同一客户端）
func (p *sqlite3Provider) Write() provider.DbClient {
	return p.client
}