package message

import "encoding/json"

// UserRecallEventMessage 用户撤回消息事件
type UserRecallEventMessage struct {
	*Message
	Event *UserRecallEventContent `json:"event"`
}

// UserRecallEventContent 用户撤回消息事件内容
type UserRecallEventContent struct {
	EventType      string `json:"event_type"`
	OpenKfID       string `json:"open_kfid"`
	ExternalUserID string `json:"external_userid"`
	RecallMsgID    string `json:"recall_msgid"`
}

var _ Messager = (*UserRecallEventMessage)(nil)

// newUserRecallEventMessage 创建用户撤回消息事件
func newUserRecallEventMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string                  `json:"msgtype"`
		Event   *UserRecallEventContent `json:"event"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &UserRecallEventMessage{
		Message: msg,
		Event:   fullMsg.Event,
	}, nil
}

// GetKind 获取消息类型
func (m *UserRecallEventMessage) GetKind() Kind {
	return MessageKindEventUserRecall
}

// Reply 回复消息
func (m *UserRecallEventMessage) Reply([]byte) ([]byte, error) {
	// 撤回事件不需要回复
	return nil, nil
}

// ServicerRecallEventMessage 接待人员撤回消息事件
type ServicerRecallEventMessage struct {
	*Message
	Event *ServicerRecallEventContent `json:"event"`
}

// ServicerRecallEventContent 接待人员撤回消息事件内容
type ServicerRecallEventContent struct {
	EventType      string `json:"event_type"`
	OpenKfID       string `json:"open_kfid"`
	ExternalUserID string `json:"external_userid"`
	RecallMsgID    string `json:"recall_msgid"`
	ServicerUserID string `json:"servicer_userid"`
}

var _ Messager = (*ServicerRecallEventMessage)(nil)

// newServicerRecallEventMessage 创建接待人员撤回消息事件
func newServicerRecallEventMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string                      `json:"msgtype"`
		Event   *ServicerRecallEventContent `json:"event"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &ServicerRecallEventMessage{
		Message: msg,
		Event:   fullMsg.Event,
	}, nil
}

// GetKind 获取消息类型
func (m *ServicerRecallEventMessage) GetKind() Kind {
	return MessageKindEventServicerRecall
}

// Reply 回复消息
func (m *ServicerRecallEventMessage) Reply([]byte) ([]byte, error) {
	// 撤回事件不需要回复
	return nil, nil
}

// RejectSwitchChangeEventMessage 拒收客户消息变更事件
type RejectSwitchChangeEventMessage struct {
	*Message
	Event *RejectSwitchChangeEventContent `json:"event"`
}

// RejectSwitchChangeEventContent 拒收消息变更事件内容
type RejectSwitchChangeEventContent struct {
	EventType      string `json:"event_type"`
	ServicerUserID string `json:"servicer_userid"`
	OpenKfID       string `json:"open_kfid"`
	ExternalUserID string `json:"external_userid"`
	RejectSwitch   uint32 `json:"reject_switch"`
}

var _ Messager = (*RejectSwitchChangeEventMessage)(nil)

// newRejectSwitchChangeEventMessage 创建拒收消息变更事件
func newRejectSwitchChangeEventMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string                          `json:"msgtype"`
		Event   *RejectSwitchChangeEventContent `json:"event"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &RejectSwitchChangeEventMessage{
		Message: msg,
		Event:   fullMsg.Event,
	}, nil
}

// GetKind 获取消息类型
func (m *RejectSwitchChangeEventMessage) GetKind() Kind {
	return MessageKindEventRejectSwitchChange
}

// Reply 回复消息
func (m *RejectSwitchChangeEventMessage) Reply([]byte) ([]byte, error) {
	// 拒收变更事件不需要回复
	return nil, nil
}
