package message

import (
	"encoding/xml"
	"time"

	"github.com/pkg/errors"
)

type VideoReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Video        *Video
}

type Video struct {
	MediaId     string
	Title       string
	Description string
}

func (m *Message) ReplyVideo(video *Video) ([]byte, error) {
	replyMsg := VideoReply{
		XMLName:      xml.Name{},
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "image",
		Video:        video,
	}

	output, err := xml.MarshalIndent(replyMsg, " ", " ")
	if err != nil {
		return nil, errors.Wrapf(err, "marshal video reply, reply: %v", replyMsg)
	}

	return output, nil
}
