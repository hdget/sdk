package wx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/pkg/errors"
)

type getAuthorizerAccessTokenRequest struct {
	ComponentAppid         string `json:"component_appid"`
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

type GetAuthorizerAccessTokenResult struct {
	api.Result
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	ExpiresIn              int    `json:"expires_in"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

const (
	// 第三方平台调用凭证 /获取授权账号调用令牌 限制：2000次/天/授权方 限制：2000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getAuthorizerAccessToken.html
	urlGetAuthorizerAccessToken = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=%s"
)

func (impl wxApiImpl) GetAuthorizerAccessToken(componentAccessToken string, authorizerAppid, authorizerRefreshToken string) (string, int, error) {
	req := &getAuthorizerAccessTokenRequest{
		ComponentAppid:         impl.appId,
		AuthorizerAppid:        authorizerAppid,
		AuthorizerRefreshToken: authorizerRefreshToken,
	}

	url := fmt.Sprintf(urlGetAuthorizerAccessToken, componentAccessToken)

	ret, err := api.Post[GetAuthorizerAccessTokenResult](url, req)
	if err != nil {
		return "", 0, errors.Wrap(err, "get authorizer access token")
	}

	if err = api.CheckResult(ret.Result, url, req); err != nil {
		return "", 0, errors.Wrap(err, "get authorizer access token")
	}

	if ret.AuthorizerAccessToken == "" {
		return "", 0, errors.New("empty authorizer access token")
	}

	return ret.AuthorizerAccessToken, ret.ExpiresIn, nil
}
