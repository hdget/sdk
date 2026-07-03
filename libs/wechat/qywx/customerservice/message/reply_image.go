package message

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ImageReplyMessage 图片回复消息
type ImageReplyMessage struct {
	ToUser  string `json:"touser"`
	ToKfID  string `json:"open_kfid"`
	MsgType string `json:"msgtype"`
	Image   struct {
		MediaID string `json:"media_id"`
	} `json:"image"`
}

// ReplyImage 回复图片消息
func (m *Message) ReplyImage(mediaID string) ([]byte, error) {
	replyMsg := ImageReplyMessage{
		ToUser:  m.ExternalUserID,
		ToKfID:  m.OpenKfID,
		MsgType: "image",
	}
	replyMsg.Image.MediaID = mediaID

	output, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal image reply, reply: %v", replyMsg)
	}

	return output, nil
}

// ReplyImageWithCode 使用code回复图片消息（事件响应）
func (m *Message) ReplyImageWithCode(code, mediaID string) ([]byte, error) {
	replyMsg := struct {
		Code    string `json:"code"`
		MsgType string `json:"msgtype"`
		Image   struct {
			MediaID string `json:"media_id"`
		} `json:"image"`
	}{
		Code:    code,
		MsgType: "image",
	}
	replyMsg.Image.MediaID = mediaID

	output, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal image reply with code, reply: %v", replyMsg)
	}

	return output, nil
}