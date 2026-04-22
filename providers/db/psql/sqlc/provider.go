package sqlc

import (
	"sync/atomic"

	"github.com/hdget/sdk/common/provider"
	"github.com/pkg/errors"
)

type sqlcProvider struct {
	defaultDb provider.DbClient
	masterDb  provider.DbClient
	slaveDbs  []provider.DbClient
	extraDbs  map[string]provider.DbClient
	slaveIdx  uint64 // 用于轮询选择 slave
}

func New(configProvider provider.Config, logger provider.Logger) (provider.Database, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	p := &sqlcProvider{
		slaveDbs: make([]provider.DbClient, len(config.Slaves)),
		extraDbs: make(map[string]provider.DbClient),
	}

	if config.Default != nil {
		p.defaultDb, err = newClient(config.Default)
		if err != nil {
			logger.Fatal("init postgresql default db connection", "err", err)
		}
		logger.Debug("init postgresql default db connection", "host", config.Default.Host)
	}

	if config.Master != nil {
		p.masterDb, err = newClient(config.Master)
		if err != nil {
			logger.Fatal("init postgresql master db connection", "err", err)
		}

		logger.Debug("init postgresql master db connection", "host", config.Master.Host)
	}

	for i, slaveConf := range config.Slaves {
		p.slaveDbs[i], err = newClient(slaveConf)
		if err != nil {
			logger.Fatal("init postgresql slave db connection", "slave", i, "err", err)
		}

		logger.Debug("init postgresql slave db connection", "index", i, "host", slaveConf.Host)
	}

	for _, extraConf := range config.Items {
		p.extraDbs[extraConf.Name], err = newClient(extraConf)
		if err != nil {
			logger.Fatal("new postgresql extra db connection", "name", extraConf.Name, "err", err)
		}

		logger.Debug("init postgresql extra db connection", "name", extraConf.Name, "host", extraConf.Host)
	}

	return p, nil
}

func NewClient(configProvider provider.Config, database ...string) (provider.DbClient, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, errors.Wrap(err, "new postgresql config")
	}

	var c *psqlConfig
	if config.Default != nil {
		c = config.Default
	} else if config.Master != nil {
		c = config.Master
	} else {
		return nil, errors.New("postgresql config not found")
	}

	// 默认使用系统默认数据库, 如果指定了数据库就用指定的
	c.Database = "postgres"
	if len(database) > 0 {
		c.Database = database[0]
	}

	client, err := newClient(c)
	if err != nil {
		return nil, errors.Wrap(err, "init postgresql sys db connection")
	}
	return client, nil
}

func (p *sqlcProvider) GetCapability() provider.Capability {
	return Capability
}

func (p *sqlcProvider) Main() provider.DbClient {
	return p.defaultDb
}

func (p *sqlcProvider) Master() provider.DbClient {
	return p.masterDb
}

func (p *sqlcProvider) Replica(i int) provider.DbClient {
	if i < 0 || i >= len(p.slaveDbs) {
		return nil
	}
	return p.slaveDbs[i]
}

func (p *sqlcProvider) Named(name string) provider.DbClient {
	return p.extraDbs[name]
}

// Read 返回用于读操作的数据库客户端（从 slave 中轮询选择，无 slave 则返回 master 或 default）
func (p *sqlcProvider) Read() provider.DbClient {
	if len(p.slaveDbs) > 0 {
		idx := atomic.AddUint64(&p.slaveIdx, 1) - 1
		return p.slaveDbs[idx%uint64(len(p.slaveDbs))]
	}
	return p.Write()
}

// Write 返回用于写操作的数据库客户端（返回 master 或 default）
func (p *sqlcProvider) Write() provider.DbClient {
	if p.masterDb != nil {
		return p.masterDb
	}
	return p.defaultDb
}