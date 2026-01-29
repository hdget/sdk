package wx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/pkg/errors"
)

// loginResult 类型
type loginResult struct {
	api.Result
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	UnionId      string `json:"unionid"`
	Scope        string `json:"scope"`
}

const (
	// 参考: https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
	urlOAuth2AccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
)

// WebAppLogin 网站应用快速扫码登录
func (impl wxApiImpl) WebAppLogin(code string) (string, string, error) {
	url := fmt.Sprintf(urlOAuth2AccessToken, impl.appId, impl.appSecret, code)

	ret, err := api.Get[loginResult](url)
	if err != nil {
		return "", "", errors.Wrap(err, "open platform web app login")
	}

	if err = api.CheckResult(ret.Result, url); err != nil {
		return "", "", errors.Wrap(err, "open platform web app login")
	}

	return ret.OpenId, ret.UnionId, nil
}
