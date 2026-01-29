package wx

import (
	"fmt"

	"github.com/elliotchance/pie/v2"
	"github.com/hdget/sdk/libs/wechat/api"
	"github.com/pkg/errors"
)

const (
	// 授权账号管理 /设置授权方选项信息 限制：1000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/setAuthorizerOptionInfo.html
	urlSetAuthorizerOption = "https://api.weixin.qq.com/cgi-bin/component/set_authorizer_option?access_token=%s"
)

func (impl wxApiImpl) SetAuthorizerOption(authorizerAccessToken string, optionName string, optionValue string) error {
	validOptionNames := []string{"location_report", "voice_recognize", "customer_service"}
	validOptionValues := []string{"0", "1"}
	switch optionName {
	case "location_report":
		validOptionValues = []string{"0", "1", "2"}
	}

	if !pie.Contains(validOptionNames, optionName) {
		return fmt.Errorf("option name not supported, optionName: %s, valid: %v", optionName, validOptionNames)
	}

	if !pie.Contains(validOptionValues, optionValue) {
		return fmt.Errorf("option value not supported, optionValue: %s, valid: %v", optionValue, validOptionValues)
	}

	req := &setAuthorizerOptionRequest{
		OptionName:  optionName,
		OptionValue: optionValue,
	}

	url := fmt.Sprintf(urlSetAuthorizerOption, authorizerAccessToken)

	ret, err := api.Post[api.Result](url, req)
	if err != nil {
		return errors.Wrap(err, "get authorizer option")
	}

	if err = api.CheckResult(ret, url, req); err != nil {
		return errors.Wrap(err, "get authorizer option")
	}

	return nil
}
