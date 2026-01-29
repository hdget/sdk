package wx

type WxAPI interface {
	Authorizer
	Credential
	Wxmp
}

type Authorizer interface {
	GetAuthorizerInfo(componentAccessToken, authorizerAppid string) (*GetAuthorizerInfoResult, error) // 授权账号管理 /获取授权账号详情
	GetAuthorizerOption(appId string, optionName string) (string, error)                              // 授权账号管理 /获取授权方选项信息
	SetAuthorizerOption(authorizerAccessToken string, optionName string, optionValue string) error    // 授权账号管理 /设置授权方选项信息
}

type Credential interface {
	StartPushComponentVerifyTicket() error                                                                                     // 第三方平台调用凭证/启动票据推送服务
	CreatePreAuthCode(componentAccessToken string) (string, int, error)                                                        // 第三方平台调用凭证/获取预授权码
	GetAuthorizerAccessToken(componentAccessToken string, authorizerAppid, authorizerRefreshToken string) (string, int, error) // 第三方平台调用凭证/获取授权账号调用令牌
	QueryAuthorizationInfo(componentAccessToken, authCode string) (*AuthorizationInfo, error)                                  // 第三方平台调用凭证/获取刷新令牌
	GetComponentAccessToken(componentVerifyTicket string) (string, int, error)                                                 // 第三方平台调用凭证 /获取令牌
}

type Wxmp interface {
	WxaCode2Session(componentAppId, componentAccessToken string, appId, code string) (*WxaCode2SessionResult, error) // 小程序登录
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
