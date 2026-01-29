package message

import (
	"fmt"
)

type LocationEventMessage struct {
	*Message
}

var (
	_ Messager = (*LocationEventMessage)(nil)
)

func newLocationEventMessage(msg *Message, data []byte) (Messager, error) {
	return &LocationEventMessage{Message: msg}, nil
}

func (m *LocationEventMessage) Reply() ([]byte, error) {
	return m.ReplyText(fmt.Sprintf("地理位置事件, message: %+v", m))
}

func (m LocationEventMessage) GetKind() Kind {
	return MessageKindEventLocation
}
