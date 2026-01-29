package wx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/pkg/errors"
)

type TemplateMessage struct {
	// 必要参数
	ToUser     string `json:"touser"`      // 发送给哪个openId
	TemplateId string `json:"template_id"` // 模板ID
	Data       map[string]struct {
		Value string `json:"value"`
	} `json:"data"` // 模板内容
	Url         string `json:"url"` // 跳转链接
	MiniProgram struct {
		AppId    string `json:"appid"`
		PagePath string `json:"pagepath"`
	} `json:"miniprogram"` // 跳小程序所需数据
	ClientMsgId string `json:"client_msg_id"` // 防重入id。对于同一个openid + client_msg_id, 只发送一条消息,10分钟有效,超过10分钟不保证效果。若无防重入需求，可不填
}

type sendTemplateMessageResult struct {
	api.Result
	Msgid int `json:"msgid"`
}

const (
	// 参考：https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
	urlSendTemplateMessage = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)

// SendTemplateMessage 发送模板消息
func (impl wxApiImpl) SendTemplateMessage(accessToken string, msg *TemplateMessage) error {
	if msg.ToUser == "" || msg.TemplateId == "" || len(msg.Data) == 0 {
		return errors.New("empty touser/template_id/data")
	}

	url := fmt.Sprintf(urlSendTemplateMessage, accessToken)

	ret, err := api.Post[sendTemplateMessageResult](url, msg)
	if err != nil {
		return errors.Wrap(err, "send template message")
	}

	if err = api.CheckResult(ret.Result, url, msg); err != nil {
		return errors.Wrap(err, "send template message")
	}

	return nil
}
