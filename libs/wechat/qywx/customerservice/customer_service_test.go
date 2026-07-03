package customerservice

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/qywx/customerservice/message"
)

// Example_usage 使用示例
func Example_usage() {
	// 1. 初始化API
	corpID := "your_corp_id"
	kfSecret := "your_kf_secret"
	api := New(corpID, kfSecret, nil)

	// 2. 校验服务器（用于企业微信回调验证）
	echostr, err := api.VerifyCorpServer("token", "signature", "timestamp", "nonce", "echostr")
	if err != nil {
		fmt.Printf("verify failed: %v\n", err)
		return
	}
	fmt.Printf("verify success: %s\n", echostr)

	// 3. 注册消息处理器
	RegisterMessageHandler(message.MessageKindText, func(m message.Messager) ([]byte, error) {
		// 处理文本消息
		textMsg := m.GetMessage()
		fmt.Printf("received text message: %s\n", textMsg.Text.Content)

		// 回复消息
		return m.Reply(nil)
	})

	// 4. 处理接收到的消息
	data := []byte(`{"msgtype":"text","text":{"content":"hello"}}`)
	reply, err := api.HandleMessage(data)
	if err != nil {
		fmt.Printf("handle message failed: %v\n", err)
		return
	}
	fmt.Printf("reply: %s\n", string(reply))

	// 5. 主动发送消息
	msg := map[string]interface{}{
		"touser":    "external_userid",
		"open_kfid": "open_kfid",
		"msgtype":   "text",
		"text":      map[string]string{"content": "您好"},
	}
	if err := api.SendMessage(msg); err != nil {
		fmt.Printf("send message failed: %v\n", err)
		return
	}
}
