package serviceaccount

import (
	"sync"

	"github.com/hdget/sdk/common/provider"
	"github.com/hdget/sdk/libs/wechat/wx"
	message2 "github.com/hdget/sdk/libs/wechat/wx/serviceaccount/message"
	"github.com/hdget/sdk/libs/wechat/wx/serviceaccount/wxapi"
)

// API 服务号
type API interface {
	wx.ApiCommon
	HandleMessage(data []byte) ([]byte, error)            // 接收普通/事件消息以及被动回复消息
	SendTemplateMessage(msg *wxapi.TemplateMessage) error // 发送模板消息
}

type serviceAccountApiImpl struct {
	wx.ApiCommon
	wxapi.Api
}

var (
	locker              sync.Mutex
	_msgKind2msgHandler = map[message2.Kind]message2.Handler{}
)

func New(appId, appSecret string, redisProvider provider.Redis) API {
	return &serviceAccountApiImpl{
		ApiCommon: wx.NewApiCommon(appId, appSecret, redisProvider),
		Api:       wxapi.New(appId, appSecret),
	}
}

// SendTemplateMessage 发送模板消息
func (impl *serviceAccountApiImpl) SendTemplateMessage(message *wxapi.TemplateMessage) error {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return err
	}
	return impl.Api.SendTemplateMessage(accessToken, message)
}

// HandleMessage 处理消息
func (impl *serviceAccountApiImpl) HandleMessage(data []byte) ([]byte, error) {
	m, err := message2.New(data)
	if err != nil {
		return nil, err
	}

	if h, exists := _msgKind2msgHandler[m.GetKind()]; exists {
		return h(m)
	}
	return m.Reply()
}

func RegisterMessageHandler(msgKind message2.Kind, handler message2.Handler) error {
	locker.Lock()
	defer locker.Unlock()
	_msgKind2msgHandler[msgKind] = handler
	return nil
}
