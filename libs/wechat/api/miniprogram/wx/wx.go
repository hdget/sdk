package wx

type WxAPI interface {
	Code2Session(code string) (*Session, error)                                               // 小程序登录凭证校验， 通过code换取openid和unionid
	CreateLimitedWxaCode(accessToken, path string, width int) ([]byte, error)                 // 创建有限小程序码
	CreateUnlimitedWxaCode(accessToken string, scene, page string, width int) ([]byte, error) // 创建无限小程序码
	GetUserPhoneNumber(accessToken, code string) (string, error)                              // 获取电话号码
}

type wxApiImpl struct {
	appId     string
	appSecret string
}

func New(appId, appSecret string) WxAPI {
	return &wxApiImpl{
		appId:     appId,
		appSecret: appSecret,
	}
}
