package sqlite3_sqlboiler

import (
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/hdget/sdk/common/types"
	"github.com/pkg/errors"
	_ "modernc.org/sqlite"
)

type sqlite3Provider struct {
	client types.DbClient
}

func New(configProvider types.ConfigProvider, logger types.LoggerProvider) (types.DbProvider, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	client, err := newClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "new sqlite3 client")
	}

	// 设置boil的缺省db
	boil.SetDB(client)

	logger.Debug("init sqlite3 provider", "db", config.DbPath)
	return &sqlite3Provider{client: client}, nil
}

// NewClient 从指定的文件创建创建数据库连接
func NewClient(dbFile string) (types.DbClient, error) {
	client, err := newClient(nil, dbFile)
	if err != nil {
		return nil, errors.Wrapf(err, "connect sqlite3: %s", dbFile)
	}

	// 设置boil的缺省db
	boil.SetDB(client)
	return client, nil
}

func (p *sqlite3Provider) GetCapability() types.Capability {
	return Capability
}

func (p *sqlite3Provider) My() types.DbClient {
	return p.client
}

func (p *sqlite3Provider) Master() types.DbClient {
	return p.client
}

func (p *sqlite3Provider) Slave(i int) types.DbClient {
	return p.client
}

func (p *sqlite3Provider) By(name string) types.DbClient {
	return p.client
}
