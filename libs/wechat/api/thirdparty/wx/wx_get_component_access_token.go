package wx

import (
	"github.com/hdget/libs/wechat/api"
	"github.com/pkg/errors"
)

type getComponentAccessTokenRequest struct {
	ComponentAppid        string `json:"component_appid"`
	ComponentAppsecret    string `json:"component_appsecret"`
	ComponentVerifyTicket string `json:"component_verify_ticket"`
}

type GetComponentAccessTokenResult struct {
	api.Result
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int    `json:"expires_in"`
}

const (
	// 第三方平台调用凭证 /获取令牌 限制：2000次/天
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getComponentAccessToken.html
	urlGetComponentAccessToken = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
)

func (impl wxApiImpl) GetComponentAccessToken(componentVerifyTicket string) (string, int, error) {
	if componentVerifyTicket == "" {
		return "", 0, errors.New("component_verify_ticket is empty, please wait at least 10 minutes")
	}

	req := &getComponentAccessTokenRequest{
		ComponentAppid:        impl.appId,
		ComponentAppsecret:    impl.appSecret,
		ComponentVerifyTicket: componentVerifyTicket,
	}

	ret, err := api.Post[GetComponentAccessTokenResult](urlGetComponentAccessToken, req)
	if err != nil {
		return "", 0, errors.Wrap(err, "get component access token")
	}

	if err = api.CheckResult(ret.Result, urlGetComponentAccessToken, req); err != nil {
		return "", 0, errors.Wrap(err, "check get component access token result")
	}

	if ret.ComponentAccessToken == "" {
		return "", 0, errors.New("empty component access token")
	}

	return ret.ComponentAccessToken, ret.ExpiresIn, nil
}
