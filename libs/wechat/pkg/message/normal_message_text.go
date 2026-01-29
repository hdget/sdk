package message

type TextNormalMessage struct {
	*Message
}

func newTextNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &TextNormalMessage{Message: msg}, nil
}

func (m *TextNormalMessage) GetKind() Kind {
	return MessageKindNormalText
}
