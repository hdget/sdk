package message

type VoiceNormalMessage struct {
	*Message
}

var (
	_ Messager = (*VoiceNormalMessage)(nil)
)

func newVoiceNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &VoiceNormalMessage{Message: msg}, nil
}

func (m *VoiceNormalMessage) GetKind() Kind {
	return MessageKindNormalVoice
}
