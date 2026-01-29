package message

import (
	"encoding/xml"
	"time"

	"github.com/pkg/errors"
)

type NewsReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	ArticleCount int
	Articles     *ArticleCollection
}

type ArticleCollection struct {
	Items []*Article `xml:"item"`
}

type Article struct {
	Title       string
	Description string
	PicUrl      string
	Url         string
}

func (m *Message) ReplyNews(articles []*Article) ([]byte, error) {
	replyMsg := &NewsReply{
		XMLName:      xml.Name{},
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "news",
		ArticleCount: len(articles),
		Articles:     &ArticleCollection{Items: articles},
	}

	output, err := xml.MarshalIndent(replyMsg, " ", " ")
	if err != nil {
		return nil, errors.Wrapf(err, "marshal news reply, reply: %v", replyMsg)
	}

	return output, nil
}
