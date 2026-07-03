package message

// Messager 消息接口
type Messager interface {
	Reply(content []byte) ([]byte, error) // 回复消息
	GetKind() Kind                        // 获取消息类型
	GetMessage() *Message                 // 获取基础消息
}

// Handler 消息处理函数类型
type Handler func(Messager) ([]byte, error)

// Message 基础消息结构
type Message struct {
	MsgType string `json:"msgtype"` // 消息类型

	// 公共字段
	OpenKfID       string `json:"open_kfid,omitempty"`       // 客服账号ID
	ExternalUserID string `json:"external_userid,omitempty"` // 客户UserID
	MsgID          int64  `json:"msgid,omitempty"`           // 消息ID
	Origin         int    `json:"origin,omitempty"`          // 消息来源

	// 事件消息公共字段（嵌套在event对象中）
	Event struct {
		EventType string `json:"event_type,omitempty"` // 事件类型
	} `json:"event,omitempty"`

	// 文本消息字段
	Text *TextContent `json:"text,omitempty"`

	// 图片消息字段
	Image *MediaContent `json:"image,omitempty"`

	// 语音消息字段
	Voice *MediaContent `json:"voice,omitempty"`

	// 视频消息字段
	Video *MediaContent `json:"video,omitempty"`

	// 文件消息字段
	File *MediaContent `json:"file,omitempty"`
}

// TextContent 文本内容
type TextContent struct {
	Content string `json:"content"`           // 文本内容
	MenuID  string `json:"menu_id,omitempty"` // 菜单ID
}

// MediaContent 媒体内容
type MediaContent struct {
	MediaID string `json:"media_id"` // 媒体文件ID
}

// GetMessage 获取基础消息
func (m *Message) GetMessage() *Message {
	return m
}

// Reply 默认回复文本消息
func (m *Message) Reply(content []byte) ([]byte, error) {
	if len(content) > 0 {
		return m.ReplyText(string(content))
	}
	// 默认回复文本消息
	return m.ReplyText("")
}
