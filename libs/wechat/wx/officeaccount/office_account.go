package officeaccount

import (
	"github.com/hdget/sdk/libs/wechat/wx"
	"github.com/hdget/sdk/libs/wechat/wx/officeaccount/wxapi"
)

type API interface {
	wx.ApiCommon
	GetJsSdkSignature(ticket, url string) (*wxapi.GetJsSdkSignatureResult, error)
	GetJsSdkTicket(accessToken string) (*wxapi.GetJsSdkTicketResult, error)
}

type officeAccountApiImpl struct {
	wx.ApiCommon
	wxapi.Api
}

func New(appId, appSecret string) API {
	return &officeAccountApiImpl{
		ApiCommon: wx.NewApiCommon(appId, appSecret),
		Api:       wxapi.New(appId, appSecret),
	}
}

func (impl officeAccountApiImpl) GetJsSdkSignature(ticket, url string) (*wxapi.GetJsSdkSignatureResult, error) {
	return impl.Api.GetJsSdkSignature(ticket, url)
}

func (impl officeAccountApiImpl) GetJsSdkTicket(accessToken string) (*wxapi.GetJsSdkTicketResult, error) {
	return impl.Api.GetJsSdkTicket(accessToken)
}
