package cache

import (
	"github.com/hdget/sdk/common/types"
	"github.com/hdget/sdk/libs/wechat/pkg/cache"
)

func SessionKey(appId string, redisProvider types.RedisProvider) cache.ObjectCache {
	return cache.NewObjectCache(appId, redisProvider, "session_key") // session key过期时间3600秒
}
