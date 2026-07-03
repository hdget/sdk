package wx

import (
	"fmt"

	"github.com/hdget/sdk/common/provider"
	"github.com/hdget/sdk/libs/wechat/pkg/api"
	"github.com/hdget/sdk/libs/wechat/pkg/cache"
	"github.com/pkg/errors"
)

type ApiCommon interface {
	GetAccessToken() (string, error)
	VerifyCallbackServer(token, signature, timestamp, nonce, echostr string) (string, error) // 校验服务号服务器
	GetAppId() string
	GetAppSecret() string
}

type apiImpl struct {
	appId         string
	appSecret     string
	redisProvider provider.Redis
}

const (
	accessTokenRedisKey = "wx:access_token"
	urlGetAccessToken   = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

func NewApiCommon(appId, appSecret string, redisProvider ...provider.Redis) ApiCommon {
	var redis provider.Redis
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
	return api.CacheFirst(
		cache.NewObjectCache(impl.appId, impl.redisProvider, accessTokenRedisKey),
		impl.wxGetAccessToken,
	)
}

// wxGetAccessToken 微信api
func (impl apiImpl) wxGetAccessToken() (string, int, error) {
	url := fmt.Sprintf(urlGetAccessToken, impl.appId, impl.appSecret)

	ret, err := api.Get[*api.AccessTokenResult](url, nil)
	if err != nil {
		return "", 0, errors.Wrapf(err, "get access token, appId: %s", impl.appId)
	}

	if err = api.CheckResult(ret.Result, url, nil); err != nil {
		return "", 0, errors.Wrapf(err, "get access token, appId: %s", impl.appId)
	}

	return ret.AccessToken, ret.ExpiresIn, nil
}
