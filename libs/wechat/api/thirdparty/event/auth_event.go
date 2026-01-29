package event

import (
	"encoding/xml"
	"strings"

	"github.com/hdget/sdk/libs/wechat/pkg/crypt"
	"github.com/pkg/errors"
)

type AuthEventKind string

const (
	AuthEventKindComponentVerifyTicket AuthEventKind = "component_verify_ticket" // 校验组件校验凭证
	AuthEventKindAuthorized            AuthEventKind = "authorized"              // 授权
	AuthEventKindUpdateAuthorized      AuthEventKind = "updateauthorized"        // 更新授权
	AuthEventKindUnauthorized          AuthEventKind = "unauthorized"            // 取消授权
)

type AuthEventHandler func(string) error // 处理事件

type authEventImpl struct {
	kind AuthEventKind
	data []byte
}

type xmlAuthEvent struct {
	InfoType string `xml:"InfoType"`
}

var (
	_authEventHandlers = map[AuthEventKind]AuthEventHandler{}
)

// NewAuthEvent 创建授权事件
func NewAuthEvent(appId, token, encodingAESKey string, message *Message) (Event, error) {
	msgCrypt, err := crypt.NewBizMsgCrypt(appId, token, encodingAESKey)
	if err != nil {
		return nil, err
	}

	data, err := msgCrypt.Decrypt(message.Signature, message.Timestamp, message.Nonce, message.Body)
	if err != nil {
		return nil, err
	}

	var evt xmlAuthEvent
	if err = xml.Unmarshal(data, &evt); err != nil {
		return nil, err
	}

	return &authEventImpl{
		kind: AuthEventKind(strings.ToLower(evt.InfoType)),
		data: data,
	}, nil
}

func RegisterAuthEventHandler(kind AuthEventKind, handler AuthEventHandler) {
	_authEventHandlers[kind] = handler
}

func (impl authEventImpl) Handle() error {
	var err error
	var processedResult string
	if p, exists := _preProcessors[impl.kind]; exists {
		processedResult, err = p.Process(impl.data)
		if err != nil {
			return errors.Wrapf(err, "pre process event, kind: %s, data: %s", impl.kind, string(impl.data))
		}
	}

	handler, exists := _authEventHandlers[impl.kind]
	if !exists {
		return errors.Wrapf(err, "handler not exists, kind: %s", impl.kind)
	}

	if err = handler(processedResult); err != nil {
		return errors.Wrapf(err, "handle authEventImpl, kind: %s, processed result: %s", impl.kind, processedResult)
	}

	return nil
}
