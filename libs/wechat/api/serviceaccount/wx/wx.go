package wx

type WxAPI interface {
	SendTemplateMessage(accessToken string, msg *TemplateMessage) error // 发送模板消息
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
