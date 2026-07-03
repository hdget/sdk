package openplatform

import (
	"github.com/hdget/sdk/libs/wechat/wx"
	"github.com/hdget/sdk/libs/wechat/wx/openplatform/wxapi"
)

type API interface {
	wx.ApiCommon
	WebAppLogin(code string) (string, string, error) // 网站应用快速扫码登录
}

type openPlatformApiImpl struct {
	wx.ApiCommon
	wxapi.Api
}

func New(appId, appSecret string) API {
	return &openPlatformApiImpl{
		ApiCommon: wx.NewApiCommon(appId, appSecret),
		Api:       wxapi.New(appId, appSecret),
	}
}

func (impl openPlatformApiImpl) WebAppLogin(code string) (string, string, error) {
	return impl.Api.WebAppLogin(code)
}
