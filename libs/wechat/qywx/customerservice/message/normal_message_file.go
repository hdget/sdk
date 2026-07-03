package message

import "encoding/json"

// FileNormalMessage 文件消息
type FileNormalMessage struct {
	*Message
}

var _ Messager = (*FileNormalMessage)(nil)

// newFileNormalMessage 创建文件消息
func newFileNormalMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType string `json:"msgtype"`
		File    struct {
			MediaID string `json:"media_id"`
		} `json:"file"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	msg.File = &MediaContent{
		MediaID: fullMsg.File.MediaID,
	}

	return &FileNormalMessage{Message: msg}, nil
}

// GetKind 获取消息类型
func (m *FileNormalMessage) GetKind() Kind {
	return MessageKindFile
}

// GetMediaID 获取媒体文件ID
func (m *FileNormalMessage) GetMediaID() string {
	if m.File == nil {
		return ""
	}
	return m.File.MediaID
}
