package message

import (
	"encoding/json"
	"fmt"

	"github.com/hdget/utils"
	"github.com/pkg/errors"
)

// New 创建消息对象
func New(data []byte) (Messager, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal msg, data: %s", utils.BytesToString(data))
	}

	// 判断是事件消息还是普通消息
	// 事件消息的 msgtype 是 "event"，并且有 event.event_type 字段
	if m.MsgType == "event" && m.Event.EventType != "" {
		return newEventMessage(&m, data)
	}

	// 根据消息类型创建具体消息
	return newNormalMessage(&m, data)
}

// newEventMessage 创建事件消息
func newEventMessage(msg *Message, data []byte) (Messager, error) {
	switch msg.Event.EventType {
	case "enter_session":
		return newEnterSessionEventMessage(msg, data)
	case "msg_send_fail":
		return newMsgSendFailEventMessage(msg, data)
	case "servicer_status_change":
		return newServicerStatusChangeEventMessage(msg, data)
	case "session_status_change":
		return newSessionStatusChangeEventMessage(msg, data)
	case "user_recall_msg":
		return newUserRecallEventMessage(msg, data)
	case "servicer_recall_msg":
		return newServicerRecallEventMessage(msg, data)
	case "reject_customer_msg_switch_change":
		return newRejectSwitchChangeEventMessage(msg, data)
	}

	return nil, fmt.Errorf("unsupported event message, event: %s", msg.Event.EventType)
}

// newNormalMessage 创建普通消息
func newNormalMessage(msg *Message, data []byte) (Messager, error) {
	switch msg.MsgType {
	case "text":
		return newTextNormalMessage(msg, data)
	case "image":
		return newImageNormalMessage(msg, data)
	case "voice":
		return newVoiceNormalMessage(msg, data)
	case "video":
		return newVideoNormalMessage(msg, data)
	case "file":
		return newFileNormalMessage(msg, data)
	case "location":
		return newLocationNormalMessage(msg, data)
	case "link":
		return newLinkNormalMessage(msg, data)
	case "business_card":
		return newBusinessCardNormalMessage(msg, data)
	case "miniprogram":
		return newMiniProgramNormalMessage(msg, data)
	case "msgmenu":
		return newMenuNormalMessage(msg, data)
	case "channels_shop_product":
		return newChannelsProductNormalMessage(msg, data)
	case "channels_shop_order":
		return newChannelsOrderNormalMessage(msg, data)
	case "merged_msg":
		return newMergedNormalMessage(msg, data)
	case "channels":
		return newChannelsNormalMessage(msg, data)
	case "meeting":
		return newMeetingNormalMessage(msg, data)
	case "schedule":
		return newScheduleNormalMessage(msg, data)
	case "note":
		return newNoteNormalMessage(msg, data)
	}

	return nil, fmt.Errorf("unsupported normal message, msgType: %s", msg.MsgType)
}