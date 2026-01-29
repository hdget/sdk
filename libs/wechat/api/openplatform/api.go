package openplatform

import (
	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/hdget/sdk/libs/wechat/api/openplatform/wx"
)

type API interface {
	api.API
	WebAppLogin(code string) (string, string, error) // 网站应用快速扫码登录
}

type openPlatformApiImpl struct {
	api.API
	wx.WxAPI
}

func New(appId, appSecret string) API {
	return &openPlatformApiImpl{
		API:   api.New(appId, appSecret),
		WxAPI: wx.New(appId, appSecret),
	}
}

func (impl openPlatformApiImpl) WebAppLogin(code string) (string, string, error) {
	return impl.WxAPI.WebAppLogin(code)
}
