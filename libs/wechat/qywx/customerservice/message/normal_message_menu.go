package message

import "encoding/json"

// MenuNormalMessage 菜单消息
type MenuNormalMessage struct {
	*Message
	MsgMenu *MenuContent `json:"msgmenu"`
}

// MenuContent 菜单内容
type MenuContent struct {
	HeadContent string     `json:"head_content"` // 起始文本
	List        []MenuItem `json:"list"`         // 菜单项配置
	TailContent string     `json:"tail_content"` // 结束文本
}

// MenuItem 菜单项
type MenuItem struct {
	Type        string           `json:"type"` // click/view/miniprogram
	Click       *MenuClick       `json:"click,omitempty"`
	View        *MenuView        `json:"view,omitempty"`
	MiniProgram *MenuMiniProgram `json:"miniprogram,omitempty"`
}

// MenuClick 点击菜单项
type MenuClick struct {
	ID      string `json:"id"`      // 菜单ID
	Content string `json:"content"` // 菜单显示内容
}

// MenuView 超链接菜单项
type MenuView struct {
	URL     string `json:"url"`     // 点击后跳转的链接
	Content string `json:"content"` // 菜单显示内容
}

// MenuMiniProgram 小程序菜单项
type MenuMiniProgram struct {
	AppID    string `json:"appid"`    // 小程序AppID
	PagePath string `json:"pagepath"` // 点击后进入的小程序页面
	Content  string `json:"content"`  // 菜单显示内容
}

var _ Messager = (*MenuNormalMessage)(nil)

// newMenuNormalMessage 创建菜单消息
func newMenuNormalMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string       `json:"msgtype"`
		MsgMenu *MenuContent `json:"msgmenu"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &MenuNormalMessage{
		Message: msg,
		MsgMenu: fullMsg.MsgMenu,
	}, nil
}

// GetKind 获取消息类型
func (m *MenuNormalMessage) GetKind() Kind {
	return MessageKindMenu
}
