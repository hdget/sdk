package message

import (
	"encoding/xml"
	"fmt"

	"github.com/hdget/utils"
	"github.com/pkg/errors"
)

func New(data []byte) (Messager, error) {
	var m Message
	err := xml.Unmarshal(data, &m)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal msg, data: %s", utils.BytesToString(data))
	}

	if m.Event != "" {
		return newEventMessage(&m, data)
	} else if m.MsgId != 0 {
		return newNormalMessage(&m, data)
	}

	return nil, fmt.Errorf("unsupported msg type, msgType: %s", m.MsgType)
}

func newEventMessage(msg *Message, data []byte) (Messager, error) {
	switch msg.Event {
	case "subscribe":
		if msg.Ticket != "" { // 未关注用户扫码
			return newUnSubscribedScanEventMessage(msg, data)
		} else { // 关注公众号
			return newSubscribeEventMessage(msg, data)
		}
	case "unsubscribe": // 取消关注公众号
		return newUnSubscribeEventMessage(msg, data)
	case "SCAN": // 已关注用户扫码
		return newSubscribedScanEventMessage(msg, data)
	case "LOCATION":
		return newLocationEventMessage(msg, data)
	case "CLICK":
		return newClickEventMessage(msg, data)
	case "VIEW":
		return newViewEventMessage(msg, data)
	}

	return nil, fmt.Errorf("unsupported event message, event: %s", msg.Event)
}

func newNormalMessage(msg *Message, data []byte) (Messager, error) {
	switch msg.MsgType {
	case "text":
		return newTextNormalMessage(msg, data)
	case "image":
		return newImageNormalMessage(msg, data)
	case "voice":
		return newVoiceNormalMessage(msg, data)
	case "video":
		return newVideoNormalMessage(msg, data)
	case "shortvideo":
		return newShortVideoNormalMessage(msg, data)
	case "location":
		return newLocationNormalMessage(msg, data)
	case "link":
		return newLinkNormalMessage(msg, data)
	}

	return nil, fmt.Errorf("unsupported normal message, event: %s", msg.MsgType)
}
