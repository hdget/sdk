package thirdparty

import (
	"fmt"
	"strings"

	"github.com/hdget/sdk/common/types"
	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/hdget/sdk/libs/wechat/api/thirdparty/cache"
	"github.com/hdget/sdk/libs/wechat/api/thirdparty/wx"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type API interface {
	api.API
	GetAuthUrl(client, redirectUrl string, authType int) (string, error) // 获取授权链接
	GetAuthorizerAppId(authCode string) (string, error)                  // 通过authCode获取授权应用的appId
	GetAuthorizerInfo(appId string) (*wx.GetAuthorizerInfoResult, error) // 获取授权应用的信息
	GetAuthorizerAccessToken(authorizerAppid string) (string, error)     // 获取授权的应用的AccessToken
	UpdateComponentVerifyTicket(componentVerifyTicket string) error      // 更新ComponentVerityTicket
}

type thirdPartyApiImpl struct {
	api.API
	wx.WxAPI
	redisProvider types.RedisProvider
}

const (
	urlPCAuth = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=%d"
	urlH5Auth = "https://open.weixin.qq.com/wxaopen/safe/bindcomponent?action=bindcomponent&no_scan=1&component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=%d#wechat_redirect"
)

func New(appId, appSecret string, redisProvider types.RedisProvider) API {
	return &thirdPartyApiImpl{
		API:           api.New(appId, appSecret, redisProvider),
		WxAPI:         wx.New(appId, appSecret),
		redisProvider: redisProvider,
	}
}

// UpdateComponentVerifyTicket 更新ComponentVerifyTicket
func (impl thirdPartyApiImpl) UpdateComponentVerifyTicket(componentVerifyTicket string) error {
	return cache.ComponentVerifyTicket(impl.GetAppId(), impl.redisProvider).Set(componentVerifyTicket)
}

// GetAuthUrl 获取授权URL
func (impl thirdPartyApiImpl) GetAuthUrl(client, redirectUrl string, authType int) (string, error) {
	componentAccessToken, err := impl.getComponentAccessToken()
	if err != nil {
		return "", err
	}

	preAuthCode, _, err := impl.WxAPI.CreatePreAuthCode(componentAccessToken)
	if err != nil {
		return "", err
	}

	// 校验authCode
	switch cast.ToInt(authType) {
	case 1, 2, 3, 4, 5, 6:
	default:
		return "", fmt.Errorf("invalid auth type: %d", authType)
	}

	// 校验redirectUrl, https://xxx
	if !strings.HasPrefix(redirectUrl, "https") {
		return "", fmt.Errorf("invalid redirect url, redirectUrl: %s", redirectUrl)
	}

	switch strings.ToLower(client) {
	case "pc":
		return fmt.Sprintf(urlPCAuth, impl.GetAppId(), preAuthCode, redirectUrl, authType), nil
	case "h5":
		return fmt.Sprintf(urlH5Auth, impl.GetAppId(), preAuthCode, redirectUrl, authType), nil
	default:
		return "", fmt.Errorf("unsupported client, client: %s", client)
	}
}

// GetAuthorizerAppId 通过授权码获取授权AppId
func (impl thirdPartyApiImpl) GetAuthorizerAppId(authCode string) (string, error) {
	if authCode == "" {
		return "", errors.New("empty auth code")
	}

	componentAccessToken, err := impl.getComponentAccessToken()
	if err != nil {
		return "", err
	}

	// 每次查询一次, accessToken可能会发生变化，需要更新缓存
	authorizationInfo, err := impl.WxAPI.QueryAuthorizationInfo(componentAccessToken, authCode)
	if err != nil {
		return "", errors.Wrap(err, "query authorization info")
	}

	// 缓存accessToken和refreshToken
	err = cache.AuthorizerAccessToken(impl.GetAppId(), impl.redisProvider, authorizationInfo.AuthorizerAppid).Set(authorizationInfo.AuthorizerAccessToken, authorizationInfo.ExpiresIn)
	if err != nil {
		return "", err
	}

	err = cache.AuthorizerRefreshToken(impl.GetAppId(), impl.redisProvider, authorizationInfo.AuthorizerAppid).Set(authorizationInfo.AuthorizerRefreshToken)
	if err != nil {
		return "", err
	}

	return authorizationInfo.AuthorizerAppid, nil
}

// GetAuthorizerInfo 获取授权应用的信息
func (impl thirdPartyApiImpl) GetAuthorizerInfo(appId string) (*wx.GetAuthorizerInfoResult, error) {
	componentAccessToken, err := impl.getComponentAccessToken()
	if err != nil {
		return nil, err
	}

	authorizerInfo, err := impl.WxAPI.GetAuthorizerInfo(componentAccessToken, appId)
	if err != nil {
		return nil, errors.Wrap(err, "get authorizer info")
	}

	return authorizerInfo, nil
}

// GetAuthorizerAccessToken 获取授权应用访问accessToken
func (impl thirdPartyApiImpl) GetAuthorizerAccessToken(authorizerAppid string) (string, error) {
	return api.CacheFirst(
		cache.AuthorizerAccessToken(impl.GetAppId(), impl.redisProvider, authorizerAppid), // cache

		func() (string, int, error) {
			componentAccessToken, err := impl.getComponentAccessToken()
			if err != nil {
				return "", 0, err
			}

			authorizerRefreshToken, err := impl.getAuthorizerRefreshToken(authorizerAppid)
			if err != nil {
				return "", 0, err
			}

			return impl.WxAPI.GetAuthorizerAccessToken(componentAccessToken, authorizerAppid, authorizerRefreshToken)
		},
	)
}

// getAuthorizerRefreshToken 获取保存的authorizerRefreshToken, 先从缓存中找，找不到从调用WX API接口获取
func (impl thirdPartyApiImpl) getAuthorizerRefreshToken(authorizerAppId string) (string, error) {
	return api.CacheFirst(
		cache.AuthorizerRefreshToken(impl.GetAppId(), impl.redisProvider, authorizerAppId), // cache

		func() (string, int, error) { // wx api
			componentAccessToken, err := impl.getComponentAccessToken()
			if err != nil {
				return "", 0, err
			}

			result, err := wx.New(impl.GetAppId(), impl.GetAppSecret()).GetAuthorizerInfo(componentAccessToken, authorizerAppId)
			if err != nil {
				return "", 0, errors.Wrapf(err, "wxapi get authorizer refresh Token, authorizerAppId: %s", authorizerAppId)
			}

			return result.Authorization.RefreshToken, 0, nil
		},
	)
}

func (impl thirdPartyApiImpl) getComponentVerifyTicket() (string, error) {
	cvtCache := cache.ComponentVerifyTicket(impl.GetAppId(), impl.redisProvider)

	componentVerifyTicket, _ := cvtCache.Get()
	if componentVerifyTicket == "" {
		// 如果缓存里面没有component verify ticket, 尝试重新推送ticket,一般要10分钟以后才会收到
		if err := wx.New(impl.GetAppId(), impl.GetAppSecret()).StartPushComponentVerifyTicket(); err != nil {
			return "", errors.Wrap(err, "start push component verify ticket")
		}
	}
	return componentVerifyTicket, nil
}

func (impl thirdPartyApiImpl) getComponentAccessToken() (string, error) {
	return api.CacheFirst(
		cache.ComponentAccessToken(impl.GetAppId(), impl.redisProvider),
		func() (string, int, error) {
			componentVerifyTicket, err := impl.getComponentVerifyTicket()
			if err != nil {
				return "", 0, err
			}

			return wx.New(impl.GetAppId(), impl.GetAppSecret()).GetComponentAccessToken(componentVerifyTicket)
		},
	)
}
