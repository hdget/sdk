package qywxapi

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/pkg/api"
	"github.com/pkg/errors"
)

// Api 企业微信客服API接口
type Api interface {
	SendMessage(accessToken string, msg interface{}) error                     // 发送消息
	ListCSR(accessToken, kfID string) ([]CSR, error)                           // 获取客服账号接待人员列表
	AddCSR(accessToken, kfID string, servicers []CSR) error                    // 添加接待人员
	DeleteCSR(accessToken, kfID string, servicers []CSR) error                 // 删除接待人员
	ListKfAccount(accessToken string) ([]KfAccount, error)                     // 获取客服账号列表
	AddKfAccount(accessToken string, account *KfAccountCreate) (string, error) // 创建客服账号
	DeleteKfAccount(accessToken, kfID string) error                            // 删除客服账号
	GetContactWay(accessToken, kfID string) (*ContactWay, error)               // 获取联系方式
	UpdateContactWay(accessToken string, contactWay *ContactWayUpdate) error   // 更新联系方式
	GetServiceState(accessToken, openKfID, externalUserID string) (int, error) // 获取会话状态
	TransServiceState(accessToken string, req *TransServiceStateReq) error     // 转换会话状态
}

type qywxApiImpl struct {
	corpID   string
	kfSecret string
}

// New 创建企业微信客服API实例
func New(corpID, kfSecret string) Api {
	return &qywxApiImpl{
		corpID:   corpID,
		kfSecret: kfSecret,
	}
}

// 消息发送相关URL
const urlSendMessage = "https://qyapi.weixin.qq.com/cgi-bin/kf/send_msg?access_token=%s"

// SendMessage 发送消息
// msg: 消息对象，可以是文本、图片、语音等消息类型
func (impl *qywxApiImpl) SendMessage(accessToken string, msg interface{}) error {
	url := fmt.Sprintf(urlSendMessage, accessToken)

	ret, err := api.Post[api.Result](url, msg)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	if err = api.CheckResult(ret, url, msg); err != nil {
		return errors.Wrap(err, "send message")
	}

	return nil
}
