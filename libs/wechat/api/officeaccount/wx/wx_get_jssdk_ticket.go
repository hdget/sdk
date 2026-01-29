package wx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/pkg/errors"
)

// GetJsSdkTicketResult 类型
type GetJsSdkTicketResult struct {
	api.Result
	Value     string `json:"ticket,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

const (
	// 参考：Js SDK签名：https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
	urlGetJsSdkTicket = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

// GetJsSdkTicket jssdk获取凭证
func (impl wxApiImpl) GetJsSdkTicket(accessToken string) (*GetJsSdkTicketResult, error) {
	url := fmt.Sprintf(urlGetJsSdkTicket, accessToken)

	ret, err := api.Get[*GetJsSdkTicketResult](url)
	if err != nil {
		return nil, errors.Wrap(err, "get js sdk ticket")
	}

	if err = api.CheckResult(ret.Result, url); err != nil {
		return nil, errors.Wrap(err, "get js sdk ticket")
	}

	return ret, nil
}
