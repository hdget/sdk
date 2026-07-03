package message

import (
	"fmt"
)

type SubscribedScanEventMessage struct {
	*Message
}

var (
	_ Messager = (*SubscribedScanEventMessage)(nil)
)

func newSubscribedScanEventMessage(msg *Message, data []byte) (Messager, error) {
	return &SubscribedScanEventMessage{Message: msg}, nil
}

func (m *SubscribedScanEventMessage) Reply() ([]byte, error) {
	return m.ReplyText(fmt.Sprintf("已关注用户扫码, message: %+v", m))
}

func (m *SubscribedScanEventMessage) GetKind() Kind {
	return MessageKindEventSubscribedScan
}
