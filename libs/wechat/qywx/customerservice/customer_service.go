package customerservice

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/hdget/sdk/common/provider"
	"github.com/hdget/sdk/libs/wechat/qywx"
	"github.com/hdget/sdk/libs/wechat/qywx/customerservice/qywxapi"
	"github.com/pkg/errors"
)

// Api 企业微信客服接口
type Api interface {
	qywx.ApiCommon
	VerifyCorpServer(token, signature, timestamp, nonce, echostr string) (string, error) // 校验企业微信服务器
	HandleMessage(data []byte) ([]byte, error)                                           // 接收消息和事件以及被动回复消息
	SendMessage(msg interface{}) error                                                   // 发送消息（主动调用）
}

type ctServiceApiImpl struct {
	qywx.ApiCommon
	qywxapi.Api
}

// New 创建企业微信客服API实例
// corpID: 企业ID
// kfSecret: 客服secret
func New(corpID, corpSecret string, redisProvider provider.Redis) Api {
	return &ctServiceApiImpl{
		ApiCommon: qywx.NewApiCommon(corpID, corpSecret, redisProvider),
		Api:       qywxapi.New(corpID, corpSecret),
	}
}

// VerifyCorpServer 企业微信接入时校验
// 参考：https://developer.work.weixin.qq.com/document/path/94699
func (impl *ctServiceApiImpl) VerifyCorpServer(token, signature, timestamp, nonce, echostr string) (string, error) {
	if signature == "" || timestamp == "" || nonce == "" {
		return "", fmt.Errorf("empty signature/timestamp/nonce")
	}

	si := []string{token, timestamp, nonce}
	sort.Strings(si)              // 字典序排序
	str := strings.Join(si, "")   // 组合字符串
	s := sha1.New()               // 返回一个新的使用SHA1校验的hash.Hash接口
	_, _ = io.WriteString(s, str) // WriteString函数将字符串数组str中的内容写入到s中
	calculatedSignature := fmt.Sprintf("%x", s.Sum(nil))

	if signature != calculatedSignature {
		return "", errors.New("signature not matched")
	}

	return echostr, nil
}
