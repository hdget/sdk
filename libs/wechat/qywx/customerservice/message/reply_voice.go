package message

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// VoiceReplyMessage 语音回复消息
type VoiceReplyMessage struct {
	ToUser  string `json:"touser"`
	ToKfID  string `json:"open_kfid"`
	MsgType string `json:"msgtype"`
	Voice   struct {
		MediaID string `json:"media_id"`
	} `json:"voice"`
}

// ReplyVoice 回复语音消息
func (m *Message) ReplyVoice(mediaID string) ([]byte, error) {
	replyMsg := VoiceReplyMessage{
		ToUser:  m.ExternalUserID,
		ToKfID:  m.OpenKfID,
		MsgType: "voice",
	}
	replyMsg.Voice.MediaID = mediaID

	output, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal voice reply, reply: %v", replyMsg)
	}

	return output, nil
}

// ReplyVoiceWithCode 使用code回复语音消息（事件响应）
func (m *Message) ReplyVoiceWithCode(code, mediaID string) ([]byte, error) {
	replyMsg := struct {
		Code    string `json:"code"`
		MsgType string `json:"msgtype"`
		Voice   struct {
			MediaID string `json:"media_id"`
		} `json:"voice"`
	}{
		Code:    code,
		MsgType: "voice",
	}
	replyMsg.Voice.MediaID = mediaID

	output, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal voice reply with code, reply: %v", replyMsg)
	}

	return output, nil
}