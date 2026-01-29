package wx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/pkg/errors"
)

// Session wechat miniprogram login session
type Session struct {
	api.Result
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

const (
	// Code2Session 小程序登录凭证校验， 通过code换取openid和unionid
	// 参考：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
	urlLogin = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func (impl wxApiImpl) Code2Session(code string) (*Session, error) {
	// 登录凭证校验
	url := fmt.Sprintf(urlLogin, impl.appId, impl.appSecret, code)

	ret, err := api.Get[*Session](url)
	if err != nil {
		return nil, errors.Wrap(err, "mini program code to session")
	}

	if err = api.CheckResult(ret.Result, url); err != nil {
		return nil, errors.Wrap(err, "mini program code to session")
	}

	return ret, nil
}
