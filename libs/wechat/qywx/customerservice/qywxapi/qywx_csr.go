package qywxapi

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/pkg/api"
	"github.com/pkg/errors"
)

// 接待人员管理相关URL
const (
	urlListCSR   = "https://qyapi.weixin.qq.com/cgi-bin/kf/servicer/list?access_token=%s&open_kfid=%s"
	urlAddCSR    = "https://qyapi.weixin.qq.com/cgi-bin/kf/servicer/add?access_token=%s"
	urlDeleteCSR = "https://qyapi.weixin.qq.com/cgi-bin/kf/servicer/del?access_token=%s"
)

// CSR 接待人员
type CSR struct {
	UserID string `json:"userid"` // 接待人员的userid
	Status int    `json:"status"` // 接待人员的接待状态: 0:未接待, 1:接待中
}

// ListCSR 获取客服账号接待人员列表
// accessToken: 访问令牌
// kfID: 客服账号ID
func (impl *qywxApiImpl) ListCSR(accessToken, kfID string) ([]CSR, error) {
	url := fmt.Sprintf(urlListCSR, accessToken, kfID)

	type servicerListResult struct {
		api.Result
		ServicerList []CSR `json:"servicer_list"`
	}

	ret, err := api.Get[servicerListResult](url)
	if err != nil {
		return nil, errors.Wrap(err, "get servicer list")
	}

	if err = api.CheckResult(ret.Result, url, nil); err != nil {
		return nil, errors.Wrap(err, "get servicer list")
	}

	return ret.ServicerList, nil
}

// AddCSR 添加接待人员
func (impl *qywxApiImpl) AddCSR(accessToken, kfID string, servicers []CSR) error {
	url := fmt.Sprintf(urlAddCSR, accessToken)

	req := struct {
		OpenKfID     string `json:"open_kfid"`
		ServicerList []CSR  `json:"servicer_list"`
	}{
		OpenKfID:     kfID,
		ServicerList: servicers,
	}

	ret, err := api.Post[api.Result](url, req)
	if err != nil {
		return errors.Wrap(err, "add servicer")
	}

	if err = api.CheckResult(ret, url, req); err != nil {
		return errors.Wrap(err, "add servicer")
	}

	return nil
}

// DeleteCSR 删除接待人员
func (impl *qywxApiImpl) DeleteCSR(accessToken, kfID string, servicers []CSR) error {
	url := fmt.Sprintf(urlDeleteCSR, accessToken)

	req := struct {
		OpenKfID     string `json:"open_kfid"`
		ServicerList []CSR  `json:"servicer_list"`
	}{
		OpenKfID:     kfID,
		ServicerList: servicers,
	}

	ret, err := api.Post[api.Result](url, req)
	if err != nil {
		return errors.Wrap(err, "delete servicer")
	}

	if err = api.CheckResult(ret, url, req); err != nil {
		return errors.Wrap(err, "delete servicer")
	}

	return nil
}
