package message

import "encoding/json"

// MsgSendFailEventMessage 消息发送失败事件
type MsgSendFailEventMessage struct {
	*Message
	Event *MsgSendFailEventContent `json:"event"`
}

// MsgSendFailEventContent 消息发送失败事件内容
type MsgSendFailEventContent struct {
	EventType      string `json:"event_type"`
	OpenKfID       string `json:"open_kfid"`
	ExternalUserID string `json:"external_userid"`
	FailMsgID      string `json:"fail_msgid"`
	FailType       uint32 `json:"fail_type"`
}

var _ Messager = (*MsgSendFailEventMessage)(nil)

// newMsgSendFailEventMessage 创建消息发送失败事件
func newMsgSendFailEventMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string                   `json:"msgtype"`
		Event   *MsgSendFailEventContent `json:"event"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &MsgSendFailEventMessage{
		Message: msg,
		Event:   fullMsg.Event,
	}, nil
}

// GetKind 获取消息类型
func (m *MsgSendFailEventMessage) GetKind() Kind {
	return MessageKindEventMsgSendFail
}

// Reply 回复消息
func (m *MsgSendFailEventMessage) Reply([]byte) ([]byte, error) {
	// 消息发送失败不需要回复
	return nil, nil
}
