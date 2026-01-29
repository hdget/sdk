package wx

import (
	"fmt"

	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/pkg/errors"
)

type UserInfoResult struct {
	api.Result
	Subscribe      int8   `json:"subscribe"`
	Openid         string `json:"openid"`
	Language       string `json:"language"`
	SubscribeTime  int64  `json:"subscribe_time"`
	UnionId        string `json:"unionid"`
	Remark         string `json:"remark"`
	GroupId        int    `json:"groupid"`
	TagIdList      []int  `json:"tagid_list"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        int    `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`
}

const (
	// 参考：https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
	urlGetUnionId = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
)

// GetUserInfo 通过openId获取用户信息
func (impl wxApiImpl) GetUserInfo(accessToken, openid string) (*UserInfoResult, error) {
	url := fmt.Sprintf(urlGetUnionId, accessToken, openid)

	ret, err := api.Get[*UserInfoResult](url)
	if err != nil {
		return nil, errors.Wrap(err, "get user info")
	}

	if err = api.CheckResult(ret.Result, url); err != nil {
		return nil, errors.Wrap(err, "get user info")
	}

	return ret, nil
}
