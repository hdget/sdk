package message

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// TextReplyMessage 文本回复消息
type TextReplyMessage struct {
	ToUser  string `json:"touser"`
	ToKfID  string `json:"open_kfid"`
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// ReplyText 回复文本消息
func (m *Message) ReplyText(content string) ([]byte, error) {
	replyMsg := TextReplyMessage{
		ToUser:  m.ExternalUserID,
		ToKfID:  m.OpenKfID,
		MsgType: "text",
	}
	replyMsg.Text.Content = content

	output, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal text reply, reply: %v", replyMsg)
	}

	return output, nil
}

// ReplyTextWithCode 使用code回复文本消息（事件响应）
func (m *Message) ReplyTextWithCode(code, content string) ([]byte, error) {
	replyMsg := struct {
		Code    string `json:"code"`
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}{
		Code:    code,
		MsgType: "text",
	}
	replyMsg.Text.Content = content

	output, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal text reply with code, reply: %v", replyMsg)
	}

	return output, nil
}