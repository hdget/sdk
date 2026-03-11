package logistics_kdniao

import (
	"github.com/hdget/sdk/common/types"
)

type kdniaoProvider struct {
	defaultClient types.LogisticsClient            // 缺省客户端
	extraClients  map[string]types.LogisticsClient // 额外客户端
}

func New(configProvider types.ConfigProvider, logger types.LoggerProvider) (types.LogisticsProvider, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	p := &kdniaoProvider{
		extraClients: make(map[string]types.LogisticsClient),
	}

	if config.Default != nil {
		p.defaultClient, err = newKdniaoClient(config.Default)
		if err != nil {
			logger.Fatal("init kdniao default client", "err", err)
		}
		logger.Debug("init kdniao default client", "ebusiness_id", config.Default.EBusinessID)
	}

	for _, itemConf := range config.Items {
		p.extraClients[itemConf.Name], err = newKdniaoClient(itemConf)
		if err != nil {
			logger.Fatal("new kdniao extra client", "name", itemConf.Name, "err", err)
		}
		logger.Debug("init kdniao extra client", "name", itemConf.Name)
	}

	return p, nil
}

func (p *kdniaoProvider) GetCapability() types.Capability {
	return Capability
}

func (p *kdniaoProvider) My() types.LogisticsClient {
	return p.defaultClient
}

func (p *kdniaoProvider) By(name string) types.LogisticsClient {
	return p.extraClients[name]
}
