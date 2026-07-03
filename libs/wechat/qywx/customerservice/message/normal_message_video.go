package message

import "encoding/json"

// VideoNormalMessage 视频消息
type VideoNormalMessage struct {
	*Message
}

var _ Messager = (*VideoNormalMessage)(nil)

// newVideoNormalMessage 创建视频消息
func newVideoNormalMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string `json:"msgtype"`
		Video   struct {
			MediaID string `json:"media_id"`
		} `json:"video"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	msg.Video = &MediaContent{
		MediaID: fullMsg.Video.MediaID,
	}

	return &VideoNormalMessage{Message: msg}, nil
}

// GetKind 获取消息类型
func (m *VideoNormalMessage) GetKind() Kind {
	return MessageKindVideo
}

// GetMediaID 获取媒体文件ID
func (m *VideoNormalMessage) GetMediaID() string {
	if m.Video == nil {
		return ""
	}
	return m.Video.MediaID
}
