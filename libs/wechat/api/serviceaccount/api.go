package serviceaccount

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/hdget/sdk/common/types"
	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/hdget/sdk/libs/wechat/api/serviceaccount/wx"
	"github.com/pkg/errors"
)

// API 服务号
type API interface {
	api.API
	VerifyServer(token, signature, timestamp, nonce, echostr string) (string, error) // 校验服务号服务器
	HandleMessage(data []byte) ([]byte, error)                                       // 接收普通/事件消息以及被动回复消息
	SendTemplateMessage(msg *wx.TemplateMessage) error                               // 发送模板消息
}

type serviceAccountApiImpl struct {
	api.API
	wx.WxAPI
}

func New(appId, appSecret string, redisProvider types.RedisProvider) API {
	return &serviceAccountApiImpl{
		API:   api.New(appId, appSecret, redisProvider),
		WxAPI: wx.New(appId, appSecret),
	}
}

// SendTemplateMessage 发送模板消息
func (impl *serviceAccountApiImpl) SendTemplateMessage(message *wx.TemplateMessage) error {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return err
	}
	return impl.WxAPI.SendTemplateMessage(accessToken, message)
}

// VerifyServer 服务号接入时校验
// 参考： https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
func (impl *serviceAccountApiImpl) VerifyServer(token, signature, timestamp, nonce, echostr string) (string, error) {
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
