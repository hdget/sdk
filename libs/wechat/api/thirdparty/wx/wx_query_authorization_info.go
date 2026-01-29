package wx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/pkg/errors"
)

type queryAuthorizationInfoRequest struct {
	ComponentAppid    string `json:"component_appid"`
	AuthorizationCode string `json:"authorization_code"`
}

type queryAuthorizationInfoResult struct {
	api.Result
	AuthorizationInfo AuthorizationInfo `json:"authorization_info"`
}

type AuthorizationInfo struct {
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	ExpiresIn              int    `json:"expires_in"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

const (
	// 第三方平台调用凭证/获取刷新令牌， 限制：2000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getAuthorizerRefreshToken.html
	urlQueryAuthorizationInfo = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=%s"
)

// QueryAuthorizationInfo 通过authCode获取AuthorizationInfo
func (impl wxApiImpl) QueryAuthorizationInfo(componentAccessToken, authCode string) (*AuthorizationInfo, error) {
	req := &queryAuthorizationInfoRequest{
		ComponentAppid:    impl.appId,
		AuthorizationCode: authCode,
	}

	url := fmt.Sprintf(urlQueryAuthorizationInfo, componentAccessToken)

	ret, err := api.Post[queryAuthorizationInfoResult](url, req)
	if err != nil {
		return nil, errors.Wrap(err, "get authorization info")
	}

	if err = api.CheckResult(ret.Result, url, req); err != nil {
		return nil, errors.Wrap(err, "get authorization info")
	}

	if ret.AuthorizationInfo.AuthorizerAccessToken == "" {
		return nil, errors.New("empty authorizer access token")
	}

	return &ret.AuthorizationInfo, nil
}
