package officeaccount

import (
	"github.com/hdget/libs/wechat/api"
	"github.com/hdget/libs/wechat/api/officeaccount/wx"
)

type API interface {
	api.API
	GetJsSdkSignature(ticket, url string) (*wx.GetJsSdkSignatureResult, error)
	GetJsSdkTicket(accessToken string) (*wx.GetJsSdkTicketResult, error)
}

type officeAccountApiImpl struct {
	api.API
	wx.WxAPI
}

func New(appId, appSecret string) API {
	return &officeAccountApiImpl{
		API:   api.New(appId, appSecret),
		WxAPI: wx.New(appId, appSecret),
	}
}

func (impl officeAccountApiImpl) GetJsSdkSignature(ticket, url string) (*wx.GetJsSdkSignatureResult, error) {
	return impl.WxAPI.GetJsSdkSignature(ticket, url)
}

func (impl officeAccountApiImpl) GetJsSdkTicket(accessToken string) (*wx.GetJsSdkTicketResult, error) {
	return impl.WxAPI.GetJsSdkTicket(accessToken)
}
