package qywxapi

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/pkg/api"
	"github.com/pkg/errors"
)

// 联系方式管理相关URL
const (
	urlGetContactWay    = "https://qyapi.weixin.qq.com/cgi-bin/kf/contact_way/get?access_token=%s"
	urlUpdateContactWay = "https://qyapi.weixin.qq.com/cgi-bin/kf/contact_way/update?access_token=%s"
)

// ContactWay 联系方式
type ContactWay struct {
	OpenKfID string `json:"open_kfid"` // 客服账号ID
	QrCode   string `json:"qrcode"`    // 联系二维码URL
}

// ContactWayUpdate 更新联系方式请求
type ContactWayUpdate struct {
	OpenKfID string `json:"open_kfid"` // 客服账号ID
	QrCode   string `json:"qrcode"`    // 联系二维码URL
}

// GetContactWay 获取联系方式
func (impl *qywxApiImpl) GetContactWay(accessToken, kfID string) (*ContactWay, error) {
	url := fmt.Sprintf(urlGetContactWay, accessToken)

	req := struct {
		OpenKfID string `json:"open_kfid"`
	}{
		OpenKfID: kfID,
	}

	type contactWayResult struct {
		api.Result
		ContactWay
	}

	ret, err := api.Post[contactWayResult](url, req)
	if err != nil {
		return nil, errors.Wrap(err, "get contact way")
	}

	if err = api.CheckResult(ret.Result, url, req); err != nil {
		return nil, errors.Wrap(err, "get contact way")
	}

	return &ret.ContactWay, nil
}

// UpdateContactWay 更新联系方式
func (impl *qywxApiImpl) UpdateContactWay(accessToken string, contactWay *ContactWayUpdate) error {
	url := fmt.Sprintf(urlUpdateContactWay, accessToken)

	ret, err := api.Post[api.Result](url, contactWay)
	if err != nil {
		return errors.Wrap(err, "update contact way")
	}

	if err = api.CheckResult(ret, url, contactWay); err != nil {
		return errors.Wrap(err, "update contact way")
	}

	return nil
}
