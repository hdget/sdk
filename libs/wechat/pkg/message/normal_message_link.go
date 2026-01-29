package message

type LinkNormalMessage struct {
	*Message
}

var (
	_ Messager = (*LinkNormalMessage)(nil)
)

func newLinkNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &LinkNormalMessage{Message: msg}, nil
}

func (m *LinkNormalMessage) GetKind() Kind {
	return MessageKindNormalLink
}
