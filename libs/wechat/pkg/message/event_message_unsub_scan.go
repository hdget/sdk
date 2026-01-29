package message

import (
	"fmt"
)

type UnSubscribedScanEventMessage struct {
	*Message
}

var (
	_ Messager = (*UnSubscribedScanEventMessage)(nil)
)

func newUnSubscribedScanEventMessage(msg *Message, data []byte) (Messager, error) {
	return &UnSubscribedScanEventMessage{Message: msg}, nil
}

func (m *UnSubscribedScanEventMessage) Reply() ([]byte, error) {
	return m.ReplyText(fmt.Sprintf("未关注用户扫码, message: %+v", m))
}

func (m *UnSubscribedScanEventMessage) GetKind() Kind {
	return MessageKindEventUnSubscribedScan
}
