package message

import "encoding/json"

// ServicerStatusChangeEventMessage 接待人员接待状态变更事件
type ServicerStatusChangeEventMessage struct {
	*Message
	Event *ServicerStatusChangeEventContent `json:"event"`
}

// ServicerStatusChangeEventContent 接待人员状态变更事件内容
type ServicerStatusChangeEventContent struct {
	EventType      string `json:"event_type"`
	ServicerUserID string `json:"servicer_userid"`
	Status         uint32 `json:"status"`
	StopType       uint32 `json:"stop_type"`
	OpenKfID       string `json:"open_kfid"`
}

var _ Messager = (*ServicerStatusChangeEventMessage)(nil)

// newServicerStatusChangeEventMessage 创建接待人员状态变更事件
func newServicerStatusChangeEventMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string                            `json:"msgtype"`
		Event   *ServicerStatusChangeEventContent `json:"event"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &ServicerStatusChangeEventMessage{
		Message: msg,
		Event:   fullMsg.Event,
	}, nil
}

// GetKind 获取消息类型
func (m *ServicerStatusChangeEventMessage) GetKind() Kind {
	return MessageKindEventServicerStatusChange
}

// Reply 回复消息
func (m *ServicerStatusChangeEventMessage) Reply([]byte) ([]byte, error) {
	// 状态变更事件不需要回复
	return nil, nil
}
