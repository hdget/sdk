package wx

import (
	"github.com/hdget/libs/wechat/api"
	"github.com/pkg/errors"
)

type startPushComponentTicketRequest struct {
	ComponentAppid     string `json:"component_appid"`
	ComponentAppSecret string `json:"component_secret"`
}

const (
	// 第三方平台调用凭证 /启动票据推送服务
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/startPushTicket.html
	urlStartPushTicket = "https://api.weixin.qq.com/cgi-bin/component/api_start_push_ticket"
)

func (impl wxApiImpl) StartPushComponentVerifyTicket() error {
	req := &startPushComponentTicketRequest{
		ComponentAppid:     impl.appId,
		ComponentAppSecret: impl.appSecret,
	}

	ret, err := api.Post[api.Result](urlStartPushTicket, req)
	if err != nil {
		return errors.Wrap(err, "start push component verify ticket")
	}

	if err = api.CheckResult(ret, urlStartPushTicket, req); err != nil {
		return errors.Wrap(err, "start push component verify ticket")
	}

	return nil
}
