package message

import "encoding/json"

// VoiceNormalMessage 语音消息
type VoiceNormalMessage struct {
	*Message
}

var _ Messager = (*VoiceNormalMessage)(nil)

// newVoiceNormalMessage 创建语音消息
func newVoiceNormalMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string `json:"msgtype"`
		Voice   struct {
			MediaID   string `json:"media_id"`
			Format    string `json:"format,omitempty"`
			Recognize string `json:"recognize,omitempty"`
			Duration  int    `json:"duration,omitempty"`
		} `json:"voice"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	msg.Voice = &MediaContent{
		MediaID: fullMsg.Voice.MediaID,
	}

	return &VoiceNormalMessage{Message: msg}, nil
}

// GetKind 获取消息类型
func (m *VoiceNormalMessage) GetKind() Kind {
	return MessageKindVoice
}

// GetMediaID 获取媒体文件ID
func (m *VoiceNormalMessage) GetMediaID() string {
	if m.Voice == nil {
		return ""
	}
	return m.Voice.MediaID
}
