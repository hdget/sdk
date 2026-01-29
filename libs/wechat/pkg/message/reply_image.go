package message

import (
	"encoding/xml"
	"time"

	"github.com/pkg/errors"
)

type ImageReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Image        *Image
}

type Image struct {
	MediaId string
}

func (m *Message) ReplyImage(img *Image) ([]byte, error) {
	replyMsg := ImageReply{
		XMLName:      xml.Name{},
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "image",
		Image:        img,
	}

	output, err := xml.MarshalIndent(replyMsg, " ", " ")
	if err != nil {
		return nil, errors.Wrapf(err, "marshal image reply, reply: %v", replyMsg)
	}

	return output, nil
}
