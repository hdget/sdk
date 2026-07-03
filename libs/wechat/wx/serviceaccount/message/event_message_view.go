package message

import (
	"fmt"
)

type ViewEventMessage struct {
	*Message
}

func newViewEventMessage(msg *Message, data []byte) (Messager, error) {
	return &ViewEventMessage{Message: msg}, nil
}

func (m *ViewEventMessage) Handle() ([]byte, error) {
	return m.ReplyText(fmt.Sprintf("跳转链接事件, message: %+v", m))
}

func (m *ViewEventMessage) GetKind() Kind {
	return MessageKindEventView
}
