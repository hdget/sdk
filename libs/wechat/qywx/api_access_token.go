package qywx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/pkg/api"
	"github.com/hdget/sdk/libs/wechat/pkg/cache"
	"github.com/pkg/errors"
)

const (
	accessTokenRedisKey = "qywx:access_token"
	urlGetAccessToken   = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

func (impl apiImpl) GetAccessToken() (string, error) {
	return api.CacheFirst(
		cache.NewObjectCache(impl.corpId, impl.redisProvider, accessTokenRedisKey),
		impl.qywxGetAccessToken,
	)
}

// qywxGetAccessToken 微信api
func (impl apiImpl) qywxGetAccessToken() (string, int, error) {
	url := fmt.Sprintf(urlGetAccessToken, impl.corpId, impl.corpSecret)

	ret, err := api.Get[*api.AccessTokenResult](url, nil)
	if err != nil {
		return "", 0, errors.Wrapf(err, "get access token, corpId: %s", impl.corpId)
	}

	if err = api.CheckResult(ret.Result, url, nil); err != nil {
		return "", 0, errors.Wrapf(err, "get access token, corpId: %s", impl.corpId)
	}

	return ret.AccessToken, ret.ExpiresIn, nil
}
