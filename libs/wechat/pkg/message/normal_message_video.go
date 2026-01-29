package message

type VideoNormalMessage struct {
	*Message
}

var (
	_ Messager = (*VideoNormalMessage)(nil)
)

func newVideoNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &VideoNormalMessage{Message: msg}, nil
}

func (m *VideoNormalMessage) GetKind() Kind {
	return MessageKindNormalVideo
}
