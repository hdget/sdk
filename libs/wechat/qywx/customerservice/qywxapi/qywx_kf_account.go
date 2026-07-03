package qywxapi

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/pkg/api"
	"github.com/pkg/errors"
)

// 客服账号管理相关URL
const (
	urlGetKfList       = "https://qyapi.weixin.qq.com/cgi-bin/kf/account/list?access_token=%s"
	urlAddKfAccount    = "https://qyapi.weixin.qq.com/cgi-bin/kf/account/add?access_token=%s"
	urlDeleteKfAccount = "https://qyapi.weixin.qq.com/cgi-bin/kf/account/del?access_token=%s"
)

// KfAccount 客服账号
type KfAccount struct {
	OpenKfID   string `json:"open_kfid"`   // 客服账号ID
	Name       string `json:"name"`        // 客服账号名称
	Avatar     string `json:"avatar"`      // 客服头像URL
	ManagePriv uint32 `json:"manage_priv"` // 管理权限
}

// KfAccountCreate 创建客服账号请求
type KfAccountCreate struct {
	Name    string `json:"name"`     // 客服账号名称
	MediaID string `json:"media_id"` // 客服头像临时素材ID
}

// ListKfAccount 获取客服账号列表
func (impl *qywxApiImpl) ListKfAccount(accessToken string) ([]KfAccount, error) {
	url := fmt.Sprintf(urlGetKfList, accessToken)

	type kfListResult struct {
		api.Result
		AccountList []KfAccount `json:"account_list"`
	}

	ret, err := api.Get[kfListResult](url)
	if err != nil {
		return nil, errors.Wrap(err, "get kf list")
	}

	if err = api.CheckResult(ret.Result, url, nil); err != nil {
		return nil, errors.Wrap(err, "get kf list")
	}

	return ret.AccountList, nil
}

// AddKfAccount 创建客服账号
func (impl *qywxApiImpl) AddKfAccount(accessToken string, account *KfAccountCreate) (string, error) {
	url := fmt.Sprintf(urlAddKfAccount, accessToken)

	type addKfResult struct {
		api.Result
		OpenKfID string `json:"open_kfid"`
	}

	ret, err := api.Post[addKfResult](url, account)
	if err != nil {
		return "", errors.Wrap(err, "add kf account")
	}

	if err = api.CheckResult(ret.Result, url, account); err != nil {
		return "", errors.Wrap(err, "add kf account")
	}

	return ret.OpenKfID, nil
}

// DeleteKfAccount 删除客服账号
func (impl *qywxApiImpl) DeleteKfAccount(accessToken, kfID string) error {
	url := fmt.Sprintf(urlDeleteKfAccount, accessToken)

	req := struct {
		OpenKfID string `json:"open_kfid"`
	}{
		OpenKfID: kfID,
	}

	ret, err := api.Post[api.Result](url, req)
	if err != nil {
		return errors.Wrap(err, "delete kf account")
	}

	if err = api.CheckResult(ret, url, req); err != nil {
		return errors.Wrap(err, "delete kf account")
	}

	return nil
}
