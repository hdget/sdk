package message

// 其他普通消息的简化实现

// ChannelsProductNormalMessage 视频号商品消息
type ChannelsProductNormalMessage struct {
	*Message
}

func newChannelsProductNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &ChannelsProductNormalMessage{Message: msg}, nil
}

func (m *ChannelsProductNormalMessage) GetKind() Kind {
	return MessageKindChannelsProduct
}

// ChannelsOrderNormalMessage 视频号订单消息
type ChannelsOrderNormalMessage struct {
	*Message
}

func newChannelsOrderNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &ChannelsOrderNormalMessage{Message: msg}, nil
}

func (m *ChannelsOrderNormalMessage) GetKind() Kind {
	return MessageKindChannelsOrder
}

// MergedNormalMessage 聊天记录消息
type MergedNormalMessage struct {
	*Message
}

func newMergedNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &MergedNormalMessage{Message: msg}, nil
}

func (m *MergedNormalMessage) GetKind() Kind {
	return MessageKindMerged
}

// ChannelsNormalMessage 视频号消息
type ChannelsNormalMessage struct {
	*Message
}

func newChannelsNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &ChannelsNormalMessage{Message: msg}, nil
}

func (m *ChannelsNormalMessage) GetKind() Kind {
	return MessageKindChannels
}

// MeetingNormalMessage 会议消息
type MeetingNormalMessage struct {
	*Message
}

func newMeetingNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &MeetingNormalMessage{Message: msg}, nil
}

func (m *MeetingNormalMessage) GetKind() Kind {
	return MessageKindMeeting
}

// ScheduleNormalMessage 日程消息
type ScheduleNormalMessage struct {
	*Message
}

func newScheduleNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &ScheduleNormalMessage{Message: msg}, nil
}

func (m *ScheduleNormalMessage) GetKind() Kind {
	return MessageKindSchedule
}

// NoteNormalMessage 笔记消息
type NoteNormalMessage struct {
	*Message
}

func newNoteNormalMessage(msg *Message, data []byte) (Messager, error) {
	return &NoteNormalMessage{Message: msg}, nil
}

func (m *NoteNormalMessage) GetKind() Kind {
	return MessageKindNote
}