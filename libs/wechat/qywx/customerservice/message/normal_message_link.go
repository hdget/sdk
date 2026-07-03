package message

import "encoding/json"

// LinkNormalMessage 链接消息
type LinkNormalMessage struct {
	*Message
	Link *LinkContent `json:"link"`
}

// LinkContent 链接内容
type LinkContent struct {
	Title  string `json:"title"`   // 标题
	Desc   string `json:"desc"`    // 描述
	URL    string `json:"url"`     // 点击后跳转的链接
	PicURL string `json:"pic_url"` // 缩略图链接
}

var _ Messager = (*LinkNormalMessage)(nil)

// newLinkNormalMessage 创建链接消息
func newLinkNormalMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string       `json:"msgtype"`
		Link    *LinkContent `json:"link"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &LinkNormalMessage{
		Message: msg,
		Link:    fullMsg.Link,
	}, nil
}

// GetKind 获取消息类型
func (m *LinkNormalMessage) GetKind() Kind {
	return MessageKindLink
}
