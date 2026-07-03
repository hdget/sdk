package qywxapi_test

import (
	"fmt"
	"testing"

	"github.com/hdget/sdk/libs/wechat/qywx/customerservice/qywxapi"
)

// TestSendMsgOnEvent 测试发送事件响应消息
func TestSendMsgOnEvent(t *testing.T) {
	// 创建API实例
	_ = qywxapi.New("your_corp_id", "your_kf_secret")

	// 示例1: 发送文本欢迎语
	// 使用TransServiceState返回的msg_code
	textMsg := qywxapi.NewEventTextMessage(
		"CALLBACK_CODE", // 从事件回调或TransServiceState获取的code
		"欢迎咨询，请问有什么可以帮助您的？",
	)

	// 注意：实际使用时需要先获取access_token
	// msgID, err := api.SendMsgOnEvent(accessToken, textMsg)

	fmt.Printf("文本消息: %+v\n", textMsg)

	// 示例2: 发送菜单消息（满意度评价）
	menuMsg := qywxapi.NewEventMenuMessage(
		"CALLBACK_CODE", // 从事件回调或TransServiceState获取的code
		&qywxapi.MenuContent{
			HeadContent: "您对本次服务是否满意呢？",
			List: []qywxapi.MenuItem{
				{
					Type: "click",
					Click: &qywxapi.MenuClickItem{
						ID:      "101",
						Content: "满意",
					},
				},
				{
					Type: "click",
					Click: &qywxapi.MenuClickItem{
						ID:      "102",
						Content: "不满意",
					},
				},
				{
					Type: "view",
					View: &qywxapi.MenuViewItem{
						URL:     "https://example.com/feedback",
						Content: "点击跳转到反馈页面",
					},
				},
				{
					Type: "miniprogram",
					MiniProgram: &qywxapi.MenuMiniProgramItem{
						AppID:    "wx1234567890abcdef",
						PagePath: "pages/feedback/index",
						Content:  "打开小程序反馈",
					},
				},
				{
					Type: "text",
					Text: &qywxapi.MenuTextItem{
						Content:    "感谢您的反馈！",
						NoNewLine: 0,
					},
				},
			},
			TailContent: "感谢您的评价",
		},
	)

	fmt.Printf("菜单消息: %+v\n", menuMsg)
}

// TestTransServiceStateAndSendMsg 测试转换会话状态并发送事件响应消息的完整流程
func TestTransServiceStateAndSendMsg(t *testing.T) {
	_ = qywxapi.New("your_corp_id", "your_kf_secret")

	// 示例：将会话转入待接入池
	// 实际使用时需要先获取access_token
	/*
		accessToken := "your_access_token"
		req := &qywxapi.TransServiceStateReq{
			OpenKfID:       "wkxxxxxxxxxxxx",
			ExternalUserID: "wmxxxxxxxxxxxx",
			ServiceState:   2, // 转入待接入池
		}

		// 转换会话状态，获取msg_code
		msgCode, err := api.TransServiceState(accessToken, req)
		if err != nil {
			t.Fatalf("转换会话状态失败: %v", err)
		}

		// 使用msg_code发送排队提示语
		textMsg := qywxapi.NewEventTextMessage(msgCode, "您好，正在为您分配客服，请稍候...")
		msgID, err := api.SendMsgOnEvent(accessToken, textMsg)
		if err != nil {
			t.Fatalf("发送消息失败: %v", err)
		}

		fmt.Printf("消息发送成功，消息ID: %s\n", msgID)
	*/
}
