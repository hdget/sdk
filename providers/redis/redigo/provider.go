package redigo

import (
	"github.com/hdget/sdk/common/provider"
)

type redigoProvider struct {
	defaultClient provider.RedisClient            // 缺省redis
	extraClients  map[string]provider.RedisClient // 额外的redis
}

func New(configProvider provider.Config, logger provider.Logger) (provider.Redis, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	p := &redigoProvider{
		extraClients: make(map[string]provider.RedisClient),
	}

	if config.Default != nil {
		p.defaultClient, err = newRedisClient(config.Default)
		if err != nil {
			logger.Fatal("init redis default client", "err", err)
		}
		logger.Debug("init redis default client", "host", config.Default.Host)
	}

	for _, itemConf := range config.Items {
		p.extraClients[itemConf.Name], err = newRedisClient(itemConf)
		if err != nil {
			logger.Fatal("new redis extra client", "name", itemConf.Name, "err", err)
		}
		logger.Debug("init redis extra client", "name", itemConf.Name, "host", itemConf.Host)
	}

	return p, nil
}

func (p *redigoProvider) GetCapability() provider.Capability {
	return Capability
}

func (p *redigoProvider) My() provider.RedisClient {
	return p.defaultClient
}

func (p *redigoProvider) By(name string) provider.RedisClient {
	return p.extraClients[name]
}
