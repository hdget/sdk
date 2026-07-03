package message

import "encoding/json"

// SessionStatusChangeEventMessage 会话状态变更事件
type SessionStatusChangeEventMessage struct {
	*Message
	Event *SessionStatusChangeEventContent `json:"event"`
}

// SessionStatusChangeEventContent 会话状态变更事件内容
type SessionStatusChangeEventContent struct {
	EventType         string `json:"event_type"`
	OpenKfID          string `json:"open_kfid"`
	ExternalUserID    string `json:"external_userid"`
	ChangeType        uint32 `json:"change_type"`
	OldServicerUserID string `json:"old_servicer_userid,omitempty"`
	NewServicerUserID string `json:"new_servicer_userid,omitempty"`
	MsgCode           string `json:"msg_code,omitempty"`
}

var _ Messager = (*SessionStatusChangeEventMessage)(nil)

// newSessionStatusChangeEventMessage 创建会话状态变更事件
func newSessionStatusChangeEventMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string                           `json:"msgtype"`
		Event   *SessionStatusChangeEventContent `json:"event"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &SessionStatusChangeEventMessage{
		Message: msg,
		Event:   fullMsg.Event,
	}, nil
}

// GetKind 获取消息类型
func (m *SessionStatusChangeEventMessage) GetKind() Kind {
	return MessageKindEventSessionStatusChange
}

// Reply 回复消息
func (m *SessionStatusChangeEventMessage) Reply([]byte) ([]byte, error) {
	if m.Event == nil {
		return nil, nil
	}

	// 根据变更类型回复
	switch m.Event.ChangeType {
	case 1: // 从接待池接入会话
		if m.Event.MsgCode != "" {
			return m.ReplyTextWithCode(m.Event.MsgCode, "您好，我来为您服务")
		}
	case 3: // 结束会话
		if m.Event.MsgCode != "" {
			return m.ReplyTextWithCode(m.Event.MsgCode, "感谢您的咨询，再见！")
		}
	}

	return nil, nil
}
