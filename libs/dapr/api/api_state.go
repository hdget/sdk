package api

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/client"
	"github.com/hdget/sdk/common/namespace"
	"github.com/hdget/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// SaveState 保存状态
func (a daprApiImpl) SaveState(ctx context.Context, storeName, key string, value interface{}) error {
	data, err := utils.ToBytes(value)
	if err != nil {
		return err
	}

	c, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer c.Close()

	err = c.SaveState(ctx, namespace.Encapsulate(storeName), key, data, nil)
	if err != nil {
		return errors.Wrapf(err, "save state, store: %s, key: %s, value: %s", storeName, key, value)
	}

	return nil
}

// GetState 获取状态
func (a daprApiImpl) GetState(ctx context.Context, storeName, key string) ([]byte, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	item, err := c.GetState(ctx, namespace.Encapsulate(storeName), key, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "get state, store: %s, key: %s", storeName, key)
	}

	return item.Value, nil
}

// GetBulkState 批量获取状态
func (a daprApiImpl) GetBulkState(ctx context.Context, storeName string, keys any) (map[string][]byte, error) {
	strKeys, err := cast.ToStringSliceE(keys)
	if err != nil {
		return nil, fmt.Errorf("invalid keys, keys: %v", keys)
	}

	c, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	items, err := c.GetBulkState(ctx, namespace.Encapsulate(storeName), strKeys, nil, 100)
	if err != nil {
		return nil, errors.Wrapf(err, "get bulk state, store: %s, keys: %s", storeName, keys)
	}

	results := make(map[string][]byte, len(items))
	for _, item := range items {
		if item.Error == "" {
			results[item.Key] = item.Value
		}
	}
	return results, nil
}

// DeleteState 删除状态
func (a daprApiImpl) DeleteState(ctx context.Context, storeName, key string) error {
	c, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer c.Close()
	err = c.DeleteState(ctx, namespace.Encapsulate(storeName), key, nil)
	if err != nil {
		return errors.Wrapf(err, "delete state, store: %s, key: %s", storeName, key)
	}

	return nil
}
