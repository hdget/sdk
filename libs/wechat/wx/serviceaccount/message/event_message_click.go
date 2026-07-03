package message

import (
	"fmt"
)

type ClickEventMessage struct {
	*Message
}

var (
	_ Messager = (*ClickEventMessage)(nil)
)

func newClickEventMessage(msg *Message, data []byte) (Messager, error) {
	return &ClickEventMessage{Message: msg}, nil
}

func (m *ClickEventMessage) Reply() ([]byte, error) {
	return m.ReplyText(fmt.Sprintf("点击事件, message: %+v", m))
}

func (m *ClickEventMessage) GetKind() Kind {
	return MessageKindEventClick
}
