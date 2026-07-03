package customerservice

import (
	"sync"

	"github.com/hdget/sdk/libs/wechat/qywx/customerservice/message"
)

var (
	locker              sync.Mutex
	_msgKind2msgHandler = map[message.Kind]message.Handler{}
)

// HandleMessage 处理消息和事件
// data: 接收到的消息数据（JSON格式）
// 返回：回复消息数据，如果有错误则返回错误
func (impl *ctServiceApiImpl) HandleMessage(data []byte) ([]byte, error) {
	// 使用工厂模式创建消息对象
	m, err := message.New(data)
	if err != nil {
		return nil, err
	}

	// 查找注册的处理器
	if h, exists := _msgKind2msgHandler[m.GetKind()]; exists {
		return h(m)
	}

	// 如果没有注册处理器，使用默认回复
	return m.Reply(nil)
}

// SendMessage 发送消息（主动调用）
// 需要先获取access_token，然后调用微信接口发送消息
func (impl *ctServiceApiImpl) SendMessage(msg interface{}) error {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return err
	}
	return impl.Api.SendMessage(accessToken, msg)
}

// RegisterMessageHandler 注册消息处理器
// msgKind: 消息类型
// handler: 处理函数
func RegisterMessageHandler(msgKind message.Kind, handler message.Handler) {
	locker.Lock()
	defer locker.Unlock()
	_msgKind2msgHandler[msgKind] = handler
}

// ClearMessageHandlers 清除所有消息处理器（主要用于测试）
func ClearMessageHandlers() {
	locker.Lock()
	defer locker.Unlock()
	_msgKind2msgHandler = map[message.Kind]message.Handler{}
}
