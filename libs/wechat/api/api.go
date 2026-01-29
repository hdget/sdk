package api

import (
	"fmt"
	"time"

	"github.com/hdget/libs/wechat/pkg/cache"
	"github.com/hdget/sdk/common/types"
	"github.com/pkg/errors"
)

type API interface {
	GetAccessToken() (string, error)
	GetAppId() string
	GetAppSecret() string
}

type apiImpl struct {
	appId         string
	appSecret     string
	redisProvider types.RedisProvider
}

const (
	networkTimeout      = 3 * time.Second
	accessTokenRedisKey = "access_token"
	urlGetAccessToken   = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

func New(appId, appSecret string, redisProvider ...types.RedisProvider) API {
	var redis types.RedisProvider
	if len(redisProvider) > 0 {
		redis = redisProvider[0]
	}

	return &apiImpl{
		appId:         appId,
		appSecret:     appSecret,
		redisProvider: redis,
	}
}

func (impl apiImpl) GetAppId() string {
	return impl.appId
}

func (impl apiImpl) GetAppSecret() string {
	return impl.appSecret
}

func (impl apiImpl) GetAccessToken() (string, error) {
	return CacheFirst(
		cache.NewObjectCache(impl.appId, impl.redisProvider, accessTokenRedisKey),
		impl.wxGetAccessToken,
	)
}

// wxGetAccessToken 微信api
func (impl apiImpl) wxGetAccessToken() (string, int, error) {
	url := fmt.Sprintf(urlGetAccessToken, impl.appId, impl.appSecret)

	ret, err := Get[*AccessTokenResult](url, nil)
	if err != nil {
		return "", 0, errors.Wrapf(err, "get access token, appId: %s", impl.appId)
	}

	if err = CheckResult(ret.Result, url, nil); err != nil {
		return "", 0, errors.Wrapf(err, "get access token, appId: %s", impl.appId)
	}

	return ret.AccessToken, ret.ExpiresIn, nil
}
