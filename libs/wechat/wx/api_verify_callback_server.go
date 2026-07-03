package wx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/pkg/crypt"
	"github.com/pkg/errors"
)

// VerifyCallbackServer 回调服务器校验
// 参考： https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
func (impl apiImpl) VerifyCallbackServer(token, signature, timestamp, nonce, echostr string) (string, error) {
	if signature == "" || timestamp == "" || nonce == "" {
		return "", fmt.Errorf("empty signature/timestamp/nonce")
	}

	calculatedSignature, err := crypt.CalculateSignature(token, timestamp, nonce, echostr)
	if err != nil {
		return "", errors.Wrap(err, "calculate signature failed")
	}

	if signature != calculatedSignature {
		return "", errors.New("signature not matched")
	}

	return echostr, nil
}
