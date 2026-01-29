package message

import (
	"encoding/xml"
	"time"

	"github.com/pkg/errors"
)

type TextReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
}

func (m *Message) ReplyText(content string) ([]byte, error) {
	replyMsg := TextReply{
		XMLName:      xml.Name{},
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}

	output, err := xml.MarshalIndent(replyMsg, " ", " ")
	if err != nil {
		return nil, errors.Wrapf(err, "marshal text msg, reply: %v", replyMsg)
	}

	return output, nil
}
