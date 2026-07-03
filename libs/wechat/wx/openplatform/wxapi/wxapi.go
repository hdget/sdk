package wxapi

type Api interface {
	WebAppLogin(code string) (string, string, error) // 网站应用快速扫码登录
}

type wxApiImpl struct {
	appId     string
	appSecret string
}

func New(appId, appSecret string) Api {
	return &wxApiImpl{
		appId:     appId,
		appSecret: appSecret,
	}
}
