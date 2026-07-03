package message

// Kind 消息类型
type Kind int

const (
	MessageKindUnknown Kind = iota

	// 普通消息类型
	MessageKindText             // 文本消息
	MessageKindImage            // 图片消息
	MessageKindVoice            // 语音消息
	MessageKindVideo            // 视频消息
	MessageKindFile             // 文件消息
	MessageKindLocation         // 位置消息
	MessageKindLink             // 链接消息
	MessageKindBusinessCard     // 名片消息
	MessageKindMiniProgram      // 小程序消息
	MessageKindMenu             // 菜单消息
	MessageKindChannelsProduct  // 视频号商品消息
	MessageKindChannelsOrder    // 视频号订单消息
	MessageKindMerged           // 聊天记录消息
	MessageKindChannels         // 视频号消息
	MessageKindMeeting          // 会议消息
	MessageKindSchedule         // 日程消息
	MessageKindNote             // 笔记消息

	// 事件消息类型
	MessageKindEventEnterSession           // 进入会话事件
	MessageKindEventMsgSendFail            // 消息发送失败事件
	MessageKindEventServicerStatusChange   // 接待人员状态变更事件
	MessageKindEventSessionStatusChange    // 会话状态变更事件
	MessageKindEventUserRecall             // 用户撤回消息事件
	MessageKindEventServicerRecall         // 接待人员撤回消息事件
	MessageKindEventRejectSwitchChange     // 拒收消息变更事件
)

// String 返回消息类型的字符串表示
func (k Kind) String() string {
	switch k {
	case MessageKindText:
		return "text"
	case MessageKindImage:
		return "image"
	case MessageKindVoice:
		return "voice"
	case MessageKindVideo:
		return "video"
	case MessageKindFile:
		return "file"
	case MessageKindLocation:
		return "location"
	case MessageKindLink:
		return "link"
	case MessageKindBusinessCard:
		return "business_card"
	case MessageKindMiniProgram:
		return "miniprogram"
	case MessageKindMenu:
		return "msgmenu"
	case MessageKindChannelsProduct:
		return "channels_shop_product"
	case MessageKindChannelsOrder:
		return "channels_shop_order"
	case MessageKindMerged:
		return "merged_msg"
	case MessageKindChannels:
		return "channels"
	case MessageKindMeeting:
		return "meeting"
	case MessageKindSchedule:
		return "schedule"
	case MessageKindNote:
		return "note"
	case MessageKindEventEnterSession:
		return "event:enter_session"
	case MessageKindEventMsgSendFail:
		return "event:msg_send_fail"
	case MessageKindEventServicerStatusChange:
		return "event:servicer_status_change"
	case MessageKindEventSessionStatusChange:
		return "event:session_status_change"
	case MessageKindEventUserRecall:
		return "event:user_recall_msg"
	case MessageKindEventServicerRecall:
		return "event:servicer_recall_msg"
	case MessageKindEventRejectSwitchChange:
		return "event:reject_customer_msg_switch_change"
	default:
		return "unknown"
	}
}