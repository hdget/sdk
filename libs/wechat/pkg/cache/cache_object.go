package cache

import (
	"strings"

	"github.com/hdget/sdk/common/types"
)

type ObjectCache interface {
	Get() (string, error)
	Set(value string, expiresIn ...int) error
}

type objectCacheImpl struct {
	Cache
	redisKey string
}

func NewObjectCache(appId string, redisProvider types.RedisProvider, redisKeys ...string) ObjectCache {
	return &objectCacheImpl{
		Cache:    newCache(appId, redisProvider),
		redisKey: strings.Join(redisKeys, ":"),
	}
}

func (s objectCacheImpl) Get() (string, error) {
	return s.Cache.Get(s.redisKey)
}

func (s objectCacheImpl) Set(value string, expiresIn ...int) error {
	return s.Cache.Set(s.redisKey, value, expiresIn...)
}
