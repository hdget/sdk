package message

import (
	"encoding/xml"
	"time"

	"github.com/pkg/errors"
)

type MusicReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Music        *Music
}

type Music struct {
	Title        string
	Description  string
	MusicUrl     string
	HQMusicUrl   string
	ThumbMediaId string
}

func (m *Message) ReplyMusic(music *Music) ([]byte, error) {
	replyMsg := MusicReply{
		XMLName:      xml.Name{},
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "image",
		Music:        music,
	}

	output, err := xml.MarshalIndent(replyMsg, " ", " ")
	if err != nil {
		return nil, errors.Wrapf(err, "marshal music reply, reply: %v", replyMsg)
	}

	return output, nil
}
