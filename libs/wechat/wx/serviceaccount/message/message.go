package message

import (
	"encoding/xml"
)

type Messager interface {
	Reply() ([]byte, error)
	GetKind() Kind
	GetMessage() *Message
}

type Handler func(Messager) ([]byte, error)

type Kind int

const (
	MessageKindUnknown               Kind = iota
	MessageKindNormalText                 // 文字消息
	MessageKindNormalImage                // 图片消息
	MessageKindNormalVoice                // 语音消息
	MessageKindNormalVideo                // 视频消息
	MessageKindNormalShortVideo           // 短视频消息
	MessageKindNormalLocation             // 地理位置消息
	MessageKindNormalLink                 // 链接消息
	MessageKindEventSubscribe             // 订阅事件
	MessageKindEventUnSubscribe           // 取消订阅事件
	MessageKindEventUnSubscribedScan      // 未关注用户扫码事件
	MessageKindEventSubscribedScan        // 关注用户扫码事件
	MessageKindEventLocation              // 位置上报事件
	MessageKindEventClick                 // 点击事件
	MessageKindEventView                  // 跳转链接事件
)

type Message struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	// 事件消息: 抽取事件消息的公共部分，避免二次解析
	Event     string  // 事件消息
	EventKey  string  // 事件消息
	Ticket    string  // 事件消息
	Latitude  float64 // 地理位置上报事件
	Longitude float64 // 地理位置上报事件
	Precision float64 // 地理位置上报事件
	// 普通消息: 抽取事件消息的公共部分，避免二次解析
	MsgId        int64  // 消息ID
	MsgDataId    string // 消息数据ID
	Idx          int    // 多图文时第几篇文章，从1开始
	Content      string // 文本消息: 文本消息内容
	MediaId      string // 图片消息|语音消息|视频消息： 媒体id，可以调用获取临时素材接口拉取数据
	Format       string // 语音消息： 语音格式，如amr，speex等
	MediaId16K   string // 语音消息： 16K采样率语音消息媒体id，可以调用获取临时素材接口拉取数据，返回16K采样率amr/speex语音
	PicUrl       string // 图片消息： 图片链接（由系统生成）
	ThumbMediaId string // 视频消息： 媒体id，可以调用多媒体文件下载接口拉取数据
}

func (m *Message) GetMessage() *Message {
	return m
}

func (m *Message) Reply() ([]byte, error) {
	return m.ReplyText("")
}
