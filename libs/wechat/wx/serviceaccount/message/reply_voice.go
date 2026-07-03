package message

import (
	"encoding/xml"
	"time"

	"github.com/pkg/errors"
)

type VoiceReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Voice        *Voice
}

type Voice struct {
	MediaId string
}

func (m *Message) ReplyVoice(voice *Voice) ([]byte, error) {
	replyMsg := VoiceReply{
		XMLName:      xml.Name{},
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "voice",
		Voice:        voice,
	}

	output, err := xml.MarshalIndent(replyMsg, " ", " ")
	if err != nil {
		return nil, errors.Wrapf(err, "marshal voice reply, reply: %v", replyMsg)
	}

	return output, nil
}
