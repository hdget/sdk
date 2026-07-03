package message

import "encoding/json"

// BusinessCardNormalMessage 名片消息
type BusinessCardNormalMessage struct {
	*Message
	BusinessCard *BusinessCardContent `json:"business_card"`
}

// BusinessCardContent 名片内容
type BusinessCardContent struct {
	UserID string `json:"userid"` // 名片UserID
}

var _ Messager = (*BusinessCardNormalMessage)(nil)

// newBusinessCardNormalMessage 创建名片消息
func newBusinessCardNormalMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType      string               `json:"msgtype"`
		BusinessCard *BusinessCardContent `json:"business_card"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &BusinessCardNormalMessage{
		Message:      msg,
		BusinessCard: fullMsg.BusinessCard,
	}, nil
}

// GetKind 获取消息类型
func (m *BusinessCardNormalMessage) GetKind() Kind {
	return MessageKindBusinessCard
}
