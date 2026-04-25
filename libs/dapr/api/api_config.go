package api

import (
	"context"

	"github.com/dapr/go-sdk/client"
	"github.com/hdget/sdk/common/namespace"
	"github.com/pkg/errors"
)

// GetConfigurationItems 获取配置项
func (a daprApiImpl) GetConfigurationItems(ctx context.Context, configStore string, keys []string) (map[string]*client.ConfigurationItem, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	items, err := c.GetConfigurationItems(ctx, namespace.Encapsulate(configStore), keys)
	if err != nil {
		return nil, errors.Wrap(err, "get configuration items")
	}

	return items, nil
}

// SubscribeConfigurationItems 订阅配置项更改
func (a daprApiImpl) SubscribeConfigurationItems(ctx context.Context, configStore string, keys []string, handler client.ConfigurationHandleFunction) (string, error) {
	c, err := client.NewClient()
	if err != nil {
		return "", errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return "", errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	subscriberId, err := c.SubscribeConfigurationItems(ctx, namespace.Encapsulate(configStore), keys, handler)
	if err != nil {
		return "", errors.Wrap(err, "subscribe configuration items update")
	}
	return subscriberId, nil
}
