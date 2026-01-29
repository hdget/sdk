package message

type UnSubscribeEventMessage struct {
	*Message
}

var (
	_ Messager = (*UnSubscribeEventMessage)(nil)
)

func newUnSubscribeEventMessage(msg *Message, data []byte) (Messager, error) {
	return &UnSubscribeEventMessage{Message: msg}, nil
}

func (m *UnSubscribeEventMessage) Reply() ([]byte, error) {
	return m.ReplyText("取消关注！")
}

func (m *UnSubscribeEventMessage) GetKind() Kind {
	return MessageKindEventUnSubscribe
}
