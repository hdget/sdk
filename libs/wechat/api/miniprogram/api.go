package miniprogram

import (
	"github.com/hdget/sdk/common/types"
	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/hdget/sdk/libs/wechat/api/miniprogram/cache"
	"github.com/hdget/sdk/libs/wechat/api/miniprogram/wx"
)

type API interface {
	api.API
	Login(code string) (string, string, error)                                         // 小程序静默登录，通过code换取UnionId
	GetUserPhoneNumber(accessToken, code string) (string, error)                       // 获取用户手机号码
	CreateLimitedWxaCode(accessToken, path string, width int) ([]byte, error)          // 生成有限的小程序码
	CreateUnlimitedWxaCode(accessToken, scene, page string, width int) ([]byte, error) // CreateUnLimitedWxaCode 生成小程序码，可接受页面参数较短，生成个数不受限
}

type miniProgramApiImpl struct {
	api.API
	wx.WxAPI
	redisProvider types.RedisProvider
}

func New(appId, appSecret string, redisProvider types.RedisProvider) API {
	return &miniProgramApiImpl{
		API:           api.New(appId, appSecret),
		WxAPI:         wx.New(appId, appSecret),
		redisProvider: redisProvider,
	}
}

func (impl miniProgramApiImpl) Login(code string) (string, string, error) {
	result, err := impl.WxAPI.Code2Session(code)
	if err != nil {
		return "", "", err
	}

	// 保存到缓存中
	err = cache.SessionKey(impl.GetAppId(), impl.redisProvider).Set(result.SessionKey)
	if err != nil {
		return "", "", err
	}

	return result.OpenId, result.UnionId, err
}

func (impl miniProgramApiImpl) GetUserPhoneNumber(accessToken, code string) (string, error) {
	return impl.WxAPI.GetUserPhoneNumber(accessToken, code)
}

func (impl miniProgramApiImpl) CreateLimitedWxaCode(accessToken, path string, width int) ([]byte, error) {
	return impl.WxAPI.CreateLimitedWxaCode(accessToken, path, width)
}

func (impl miniProgramApiImpl) CreateUnlimitedWxaCode(accessToken, scene, page string, width int) ([]byte, error) {
	return impl.WxAPI.CreateUnlimitedWxaCode(accessToken, scene, page, width)
}
