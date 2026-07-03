package qywxapi

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/pkg/api"
	"github.com/pkg/errors"
)

// 事件响应消息相关URL
const urlSendMsgOnEvent = "https://qyapi.weixin.qq.com/cgi-bin/kf/send_msg_on_event?access_token=%s"

// EventTextMessage 事件响应文本消息
type EventTextMessage struct {
	Code    string         `json:"code"`    // 事件响应消息对应的code，通过事件回调下发，仅可使用一次
	MsgID   string         `json:"msgid"`   // 消息ID，如果请求参数指定了msgid，则原样返回，否则系统自动生成并返回。不多于32字节
	MsgType string         `json:"msgtype"` // 消息类型，此时固定为：text
	Text    *TextContent   `json:"text"`    // 文本消息内容
}

// EventMenuMessage 事件响应菜单消息
type EventMenuMessage struct {
	Code     string        `json:"code"`              // 事件响应消息对应的code，通过事件回调下发，仅可使用一次
	MsgID    string        `json:"msgid,omitempty"`   // 消息ID，如果请求参数指定了msgid，则原样返回，否则系统自动生成并返回。不多于32字节
	MsgType  string        `json:"msgtype"`           // 消息类型，此时固定为：msgmenu
	MsgMenu  *MenuContent  `json:"msgmenu"`           // 菜单消息内容
}

// TextContent 文本内容
type TextContent struct {
	Content string `json:"content"` // 消息内容，最长不超过2048个字节
}

// MenuContent 菜单消息内容
type MenuContent struct {
	HeadContent string       `json:"head_content,omitempty"` // 起始文本，不多于1024字节
	List        []MenuItem   `json:"list,omitempty"`         // 菜单项配置，不超过10个
	TailContent string       `json:"tail_content,omitempty"` // 结束文本，不多于1024字节
}

// MenuItem 菜单项
type MenuItem struct {
	Type       string           `json:"type"`                 // 菜单类型：click-回复菜单，view-超链接菜单，miniprogram-小程序菜单，text-文本
	Click      *MenuClickItem   `json:"click,omitempty"`      // type为click的菜单项
	View       *MenuViewItem    `json:"view,omitempty"`       // type为view的菜单项
	MiniProgram *MenuMiniProgramItem `json:"miniprogram,omitempty"` // type为miniprogram的菜单项
	Text       *MenuTextItem    `json:"text,omitempty"`       // type为text的菜单项
}

// MenuClickItem 回复菜单项
type MenuClickItem struct {
	ID      string `json:"id"`                // 菜单ID。不少于1字节，不多于128字节
	Content string `json:"content,omitempty"` // 菜单显示内容。不少于1字节，不多于128字节
}

// MenuViewItem 超链接菜单项
type MenuViewItem struct {
	URL     string `json:"url"`               // 点击后跳转的链接。不少于1字节，不多于2048字节
	Content string `json:"content,omitempty"` // 菜单显示内容。不少于1字节，不多于1024字节
}

// MenuMiniProgramItem 小程序菜单项
type MenuMiniProgramItem struct {
	AppID   string `json:"appid"`               // 小程序appid。不少于1字节，不多于32字节
	PagePath string `json:"pagepath"`           // 点击后进入的小程序页面。不少于1字节，不多于1024字节
	Content string `json:"content,omitempty"`   // 菜单显示内容。不多于1024字节
}

// MenuTextItem 文本菜单项
type MenuTextItem struct {
	Content    string `json:"content"`              // 文本内容，支持\n换行。不少于1字节，不多于256字节
	NoNewLine int    `json:"no_newline,omitempty"` // 内容后面是否不换行，0-换行 1-不换行，默认为0
}

// NewEventTextMessage 创建事件响应文本消息
func NewEventTextMessage(code, content string) *EventTextMessage {
	return &EventTextMessage{
		Code:    code,
		MsgType: "text",
		Text: &TextContent{
			Content: content,
		},
	}
}

// NewEventMenuMessage 创建事件响应菜单消息
func NewEventMenuMessage(code string, menu *MenuContent) *EventMenuMessage {
	return &EventMenuMessage{
		Code:    code,
		MsgType: "msgmenu",
		MsgMenu: menu,
	}
}

// SendMsgOnEvent 发送事件响应消息
// msg: 消息对象，可以是EventTextMessage或EventMenuMessage
func (impl *qywxApiImpl) SendMsgOnEvent(accessToken string, msg interface{}) (string, error) {
	url := fmt.Sprintf(urlSendMsgOnEvent, accessToken)

	type sendMsgOnEventResult struct {
		api.Result
		MsgID string `json:"msgid"` // 消息ID
	}

	ret, err := api.Post[sendMsgOnEventResult](url, msg)
	if err != nil {
		return "", errors.Wrap(err, "send msg on event")
	}

	if err = api.CheckResult(ret.Result, url, msg); err != nil {
		return "", errors.Wrap(err, "send msg on event")
	}

	return ret.MsgID, nil
}
