package cache

import (
	"github.com/hdget/libs/wechat/pkg/cache"
	"github.com/hdget/sdk/common/types"
)

func SessionKey(appId string, redisProvider types.RedisProvider) cache.ObjectCache {
	return cache.NewObjectCache(appId, redisProvider, "session_key") // session key过期时间3600秒
}
