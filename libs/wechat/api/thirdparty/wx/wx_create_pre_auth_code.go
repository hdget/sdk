package wx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/pkg/errors"
)

type createPreAuthCodeRequest struct {
	ComponentAppid string `json:"component_appid"`
}

type createPreAuthCodeResult struct {
	api.Result
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

const (
	// 第三方平台调用凭证/获取预授权码 限制：2000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getPreAuthCode.html
	urlCreatePreAuthCode = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=%s"
)

func (impl wxApiImpl) CreatePreAuthCode(componentAccessToken string) (string, int, error) {
	req := &createPreAuthCodeRequest{
		ComponentAppid: impl.appId,
	}

	url := fmt.Sprintf(urlCreatePreAuthCode, componentAccessToken)

	ret, err := api.Post[createPreAuthCodeResult](url, req)
	if err != nil {
		return "", 0, errors.Wrap(err, "get authorizer option")
	}

	if err = api.CheckResult(ret.Result, url, req); err != nil {
		return "", 0, errors.Wrap(err, "get authorizer option")
	}

	if ret.PreAuthCode == "" {
		return "", 0, errors.New("empty pre auth code")
	}

	return ret.PreAuthCode, ret.ExpiresIn, nil
}
