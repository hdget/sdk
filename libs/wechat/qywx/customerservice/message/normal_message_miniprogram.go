package message

import "encoding/json"

// MiniProgramNormalMessage 小程序消息
type MiniProgramNormalMessage struct {
	*Message
	MiniProgram *MiniProgramContent `json:"miniprogram"`
}

// MiniProgramContent 小程序内容
type MiniProgramContent struct {
	Title        string `json:"title"`          // 标题
	AppID        string `json:"appid"`          // 小程序AppID
	PagePath     string `json:"pagepath"`       // 小程序页面路径
	ThumbMediaID string `json:"thumb_media_id"` // 小程序消息封面MediaID
}

var _ Messager = (*MiniProgramNormalMessage)(nil)

// newMiniProgramNormalMessage 创建小程序消息
func newMiniProgramNormalMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType     string              `json:"msgtype"`
		MiniProgram *MiniProgramContent `json:"miniprogram"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &MiniProgramNormalMessage{
		Message:     msg,
		MiniProgram: fullMsg.MiniProgram,
	}, nil
}

// GetKind 获取消息类型
func (m *MiniProgramNormalMessage) GetKind() Kind {
	return MessageKindMiniProgram
}
