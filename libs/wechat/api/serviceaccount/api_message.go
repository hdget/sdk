package serviceaccount

import (
	"sync"

	"github.com/hdget/sdk/libs/wechat/pkg/message"
)

var (
	locker              sync.Mutex
	_msgKind2msgHandler = map[message.Kind]message.Handler{}
)

// HandleMessage 处理消息
func (impl *serviceAccountApiImpl) HandleMessage(data []byte) ([]byte, error) {
	m, err := message.New(data)
	if err != nil {
		return nil, err
	}

	if h, exists := _msgKind2msgHandler[m.GetKind()]; exists {
		return h(m)
	}
	return m.Reply()
}

func RegisterMessageHandler(msgKind message.Kind, handler message.Handler) error {
	locker.Lock()
	defer locker.Unlock()
	_msgKind2msgHandler[msgKind] = handler
	return nil
}
