package message

type SubscribeEventMessage struct {
	*Message
}

var (
	_ Messager = (*SubscribeEventMessage)(nil)
)

func newSubscribeEventMessage(msg *Message, data []byte) (Messager, error) {
	return &SubscribeEventMessage{Message: msg}, nil
}

func (m *SubscribeEventMessage) Reply() ([]byte, error) {
	return m.ReplyText("欢迎关注！")
}

func (m *SubscribeEventMessage) GetKind() Kind {
	return MessageKindEventSubscribe
}
