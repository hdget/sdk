package message

type ShortVideoNormalMessage struct {
	*Message
}

var (
	_ Messager = (*ShortVideoNormalMessage)(nil)
)

func newShortVideoNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &ShortVideoNormalMessage{Message: msg}, nil
}

func (m *ShortVideoNormalMessage) GetKind() Kind {
	return MessageKindNormalShortVideo
}
