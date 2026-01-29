package cache

import (
	"github.com/hdget/libs/wechat/pkg/cache"
	"github.com/hdget/sdk/common/types"
)

func AuthorizerAccessToken(appId string, redisProvider types.RedisProvider, args ...string) cache.ObjectCache {
	return cache.NewObjectCache(appId, redisProvider, append([]string{"authorizer_access_token"}, args...)...)
}

func AuthorizerRefreshToken(appId string, redisProvider types.RedisProvider, args ...string) cache.ObjectCache {
	return cache.NewObjectCache(appId, redisProvider, append([]string{"authorizer_refresh_token"}, args...)...)
}

func ComponentAccessToken(appId string, redisProvider types.RedisProvider) cache.ObjectCache {
	return cache.NewObjectCache(appId, redisProvider, "component_access_token")
}

func ComponentVerifyTicket(appId string, redisProvider types.RedisProvider) cache.ObjectCache {
	return cache.NewObjectCache(appId, redisProvider, "component_verify_ticket")
}
