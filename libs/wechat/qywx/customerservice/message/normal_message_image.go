package message

import "encoding/json"

// ImageNormalMessage 图片消息
type ImageNormalMessage struct {
	*Message
}

var _ Messager = (*ImageNormalMessage)(nil)

// newImageNormalMessage 创建图片消息
func newImageNormalMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string `json:"msgtype"`
		Image   struct {
			MediaID string `json:"media_id"`
		} `json:"image"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	msg.Image = &MediaContent{
		MediaID: fullMsg.Image.MediaID,
	}

	return &ImageNormalMessage{Message: msg}, nil
}

// GetKind 获取消息类型
func (m *ImageNormalMessage) GetKind() Kind {
	return MessageKindImage
}

// GetMediaID 获取媒体文件ID
func (m *ImageNormalMessage) GetMediaID() string {
	if m.Image == nil {
		return ""
	}
	return m.Image.MediaID
}
