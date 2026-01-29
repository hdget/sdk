package message

type ImageNormalMessage struct {
	*Message
}

var (
	_ Messager = (*ImageNormalMessage)(nil)
)

func newImageNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &ImageNormalMessage{Message: msg}, nil
}

func (m *ImageNormalMessage) GetKind() Kind {
	return MessageKindNormalImage
}
