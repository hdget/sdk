package message

import "encoding/json"

// EnterSessionEventMessage 进入会话事件
type EnterSessionEventMessage struct {
	*Message
	Event *EnterSessionEventContent `json:"event"`
}

// EnterSessionEventContent 进入会话事件内容
type EnterSessionEventContent struct {
	EventType      string              `json:"event_type"`
	OpenKfID       string              `json:"open_kfid"`
	ExternalUserID string              `json:"external_userid"`
	Scene          string              `json:"scene"`
	SceneParam     string              `json:"scene_param"`
	WelcomeCode    string              `json:"welcome_code"`
	WechatChannels *WechatChannelsInfo `json:"wechat_channels,omitempty"`
}

// WechatChannelsInfo 视频号信息
type WechatChannelsInfo struct {
	Nickname     string `json:"nickname,omitempty"`
	ShopNickname string `json:"shop_nickname,omitempty"`
	Scene        uint32 `json:"scene"`
}

var _ Messager = (*EnterSessionEventMessage)(nil)

// newEnterSessionEventMessage 创建进入会话事件
func newEnterSessionEventMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string                    `json:"msgtype"`
		Event   *EnterSessionEventContent `json:"event"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &EnterSessionEventMessage{
		Message: msg,
		Event:   fullMsg.Event,
	}, nil
}

// GetKind 获取消息类型
func (m *EnterSessionEventMessage) GetKind() Kind {
	return MessageKindEventEnterSession
}

// Reply 回复消息
func (m *EnterSessionEventMessage) Reply([]byte) ([]byte, error) {
	// 如果有 welcome_code，可以发送欢迎语
	if m.Event != nil && m.Event.WelcomeCode != "" {
		return m.ReplyTextWithCode(m.Event.WelcomeCode, "您好，欢迎咨询！")
	}
	return nil, nil
}
