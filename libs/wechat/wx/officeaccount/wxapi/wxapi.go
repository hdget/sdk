package wxapi

type Api interface {
	GetJsSdkSignature(ticket, url string) (*GetJsSdkSignatureResult, error) // 生成微信签名
	GetJsSdkTicket(accessToken string) (*GetJsSdkTicketResult, error)       // 获取凭证
	GetUserInfo(accessToken, openid string) (*UserInfoResult, error)        // 通过openId获取用户信息
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
