package cache

import (
	"fmt"

	"github.com/hdget/sdk/common/types"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string, expires ...int) error
	Del(key string) error
	HGet(key, member string) (string, error)
	HSet(key, member, value string) error
}

type cacheImpl struct {
	AppId         string
	RedisProvider types.RedisProvider
}

func newCache(appId string, redisProvider types.RedisProvider) Cache {
	return &cacheImpl{
		AppId:         appId,
		RedisProvider: redisProvider,
	}
}

const (
	redisKeyTemplate = "wx:%s:%s" // wx:appid:key
)

func (c *cacheImpl) Get(key string) (string, error) {
	bs, err := c.RedisProvider.My().Get(c.getFullKey(key))
	return string(bs), err
}

func (c *cacheImpl) Set(key, value string, expires ...int) error {
	if len(expires) == 0 {
		return c.RedisProvider.My().Set(c.getFullKey(key), value)
	}
	return c.RedisProvider.My().SetEx(c.getFullKey(key), value, expires[0])
}

func (c *cacheImpl) Del(key string) error {
	return c.RedisProvider.My().Del(c.getFullKey(key))
}

func (c *cacheImpl) HGet(key, member string) (string, error) {
	return c.RedisProvider.My().HGetString(c.getFullKey(key), member)
}

func (c *cacheImpl) HSet(key, member, value string) error {
	_, err := c.RedisProvider.My().HSet(c.getFullKey(key), member, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *cacheImpl) getFullKey(key string) string {
	return fmt.Sprintf(redisKeyTemplate, c.AppId, key)
}
