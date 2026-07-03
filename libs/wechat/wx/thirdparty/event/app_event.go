package event

import (
	"encoding/xml"

	"github.com/hdget/sdk/libs/wechat/pkg/crypt"
	"github.com/hdget/utils"
)

type AppEventKind string

type AppEventHandler func() error

type appEventImpl struct {
	kind     AppEventKind
	data     []byte
	handlers map[AppEventKind]AppEventHandler
}

type xmlAppEvent struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
}

var (
	_appEventHandlers = map[AppEventKind]AppEventHandler{}
)

func NewAppEvent(appId, token, encodingAESKey string, message *Message) (Event, error) {
	msgCrypt, err := crypt.NewBizMsgCrypt(appId, token, encodingAESKey)
	if err != nil {
		return nil, err
	}

	encryptedData, err := getEncryptField(message.Body)
	if err != nil {
		return nil, err
	}

	data, err := msgCrypt.Decrypt(message.Signature, message.Timestamp, message.Nonce, encryptedData)
	if err != nil {
		return nil, err
	}

	var evt xmlAppEvent
	if err = xml.Unmarshal(data, &evt); err != nil {
		return nil, err
	}

	return &appEventImpl{
		data:     data,
		handlers: make(map[AppEventKind]AppEventHandler),
	}, nil
}

// RegisterAppEventHandler 注册代运营APP事件处理Handler
func RegisterAppEventHandler(kind AppEventKind, handler AppEventHandler) {
	_appEventHandlers[kind] = handler
}

// getEncryptField 从XML中提取加密消息字段
func getEncryptField(body string) (string, error) {
	type XmlData struct {
		Encrypt string `xml:"Encrypt"`
	}

	var data XmlData
	err := xml.Unmarshal(utils.StringToBytes(body), &data)
	if err != nil {
		return "", err
	}

	return data.Encrypt, nil
}

func (impl appEventImpl) Handle() error {
	if handler, ok := _appEventHandlers[impl.kind]; ok {
		if err := handler(); err != nil {
			return err
		}
	}

	return nil
}
