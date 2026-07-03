package qywxapi

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/pkg/api"
	"github.com/pkg/errors"
)

// 会话状态管理相关URL
const (
	urlGetServiceState   = "https://qyapi.weixin.qq.com/cgi-bin/kf/service_state/get?access_token=%s"
	urlTransServiceState = "https://qyapi.weixin.qq.com/cgi-bin/kf/service_state/trans?access_token=%s"
)

// TransServiceStateReq 转换会话状态请求
type TransServiceStateReq struct {
	OpenKfID       string `json:"open_kfid"`       // 客服账号ID
	ExternalUserID string `json:"external_userid"` // 客户UserID
	ServiceState   int    `json:"service_state"`   // 会话状态
	ServicerUserID string `json:"servicer_userid"` // 接待人员userid
}

// GetServiceState 获取会话状态
func (impl *qywxApiImpl) GetServiceState(accessToken, openKfID, externalUserID string) (int, error) {
	url := fmt.Sprintf(urlGetServiceState, accessToken)

	req := struct {
		OpenKfID       string `json:"open_kfid"`
		ExternalUserID string `json:"external_userid"`
	}{
		OpenKfID:       openKfID,
		ExternalUserID: externalUserID,
	}

	type serviceStateResult struct {
		api.Result
		ServiceState int `json:"service_state"`
	}

	ret, err := api.Post[serviceStateResult](url, req)
	if err != nil {
		return 0, errors.Wrap(err, "get service state")
	}

	if err = api.CheckResult(ret.Result, url, req); err != nil {
		return 0, errors.Wrap(err, "get service state")
	}

	return ret.ServiceState, nil
}

// TransServiceState 转换会话状态
func (impl *qywxApiImpl) TransServiceState(accessToken string, req *TransServiceStateReq) (string, error) {
	url := fmt.Sprintf(urlTransServiceState, accessToken)

	type transServiceStateResult struct {
		api.Result
		MsgCode string `json:"msg_code"` // 用于发送响应事件消息的code
	}

	ret, err := api.Post[transServiceStateResult](url, req)
	if err != nil {
		return "", errors.Wrap(err, "trans service state")
	}

	if err = api.CheckResult(ret.Result, url, req); err != nil {
		return "", errors.Wrap(err, "trans service state")
	}

	return ret.MsgCode, nil
}
