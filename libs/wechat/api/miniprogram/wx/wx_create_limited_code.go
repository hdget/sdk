package wx

import (
	"encoding/json"
	"fmt"

	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/hdget/utils/cmp"
	"github.com/pkg/errors"
)

type createLimitedWxaCodeRequest struct {
	// 要打开的小程序版本。正式版为 release，体验版为 trial，开发版为 develop
	EnvVersion string `json:"env_version"`
	// 二维码的宽度，单位 px。最小 280px，最大 1280px
	Width int `json:"width"`
	// auto_color 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调
	AutoColor bool `json:"auto_color"`
	// auto_color 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"} 十进制表示
	LineColor struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
	} `json:"line_color"`
	// 是否需要透明底色，为 true 时，生成透明底色的小程序码
	IsHyaline bool `json:"is_hyaline"`
	// 扫码进入的小程序页面路径，最大长度 128 字节，不能为空；
	// 对于小游戏，可以只传入 query 部分，来实现传参效果，如：传入 "?foo=bar"，
	// 即可在 wx.getLaunchOptionsSync 接口中的 query 参数获取到 {foo:"bar"}。
	Path string `json:"path"`
}

const (
	// 小程序码与小程序链接 /小程序码 /获取小程序码，可接受path参数较长，生成个数受限
	// 参考：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/createQRCode.html
	urlGetLimitedWxaCode = "https://api.weixin.qq.com/wxa/getwxacode?access_token=%s"
)

// CreateLimitedWxaCode 创建小程序码
func (impl wxApiImpl) CreateLimitedWxaCode(accessToken, path string, width int) ([]byte, error) {
	// 获取post的内容
	req := &createLimitedWxaCodeRequest{
		Path:       path,
		EnvVersion: "release",
		Width:      width,
		AutoColor:  true,
	}

	url := fmt.Sprintf(urlGetLimitedWxaCode, accessToken)

	data, err := api.PostResponse(url, req)
	if err != nil {
		return nil, errors.Wrap(err, "get wxa code")
	}

	// 如果不是图像数据，那就是json错误数据
	if !cmp.IsImageData(data) {
		var result api.Result
		err = json.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(result.ErrMsg)
	}

	return data, nil
}
