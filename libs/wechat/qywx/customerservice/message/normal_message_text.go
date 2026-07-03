package message

import "encoding/json"

// TextNormalMessage 文本消息
type TextNormalMessage struct {
	*Message
}

var _ Messager = (*TextNormalMessage)(nil)

// newTextNormalMessage 创建文本消息
func newTextNormalMessage(msg *Message, data []byte) (Messager, error) {
	// 重新解析完整的文本消息结构
	var fullMsg struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
			MenuID  string `json:"menu_id,omitempty"`
		} `json:"text"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	msg.Text = &TextContent{
		Content: fullMsg.Text.Content,
		MenuID:  fullMsg.Text.MenuID,
	}

	return &TextNormalMessage{Message: msg}, nil
}

// GetKind 获取消息类型
func (m *TextNormalMessage) GetKind() Kind {
	return MessageKindText
}
