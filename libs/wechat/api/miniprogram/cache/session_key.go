package cache

import (
	"github.com/hdget/sdk/common/provider"
	"github.com/hdget/sdk/libs/wechat/pkg/cache"
)

func SessionKey(appId string, redisProvider provider.Redis) cache.ObjectCache {
	return cache.NewObjectCache(appId, redisProvider, "session_key") // session key过期时间3600秒
}
