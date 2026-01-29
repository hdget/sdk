package message

type LocationNormalMessage struct {
	*Message
}

var (
	_ Messager = (*LocationNormalMessage)(nil)
)

func newLocationNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &LocationNormalMessage{Message: msg}, nil
}

func (m *LocationNormalMessage) GetKind() Kind {
	return MessageKindNormalLocation
}
