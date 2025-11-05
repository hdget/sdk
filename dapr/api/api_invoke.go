package api

import (
	"encoding/json"

	"github.com/dapr/go-sdk/client"
	"github.com/hdget/common/namespace"
	"github.com/hdget/sdk/dapr/utils"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

const ContentTypeJson = "application/json"

// Invoke 调用dapr服务
func (a apiImpl) Invoke(app string, apiVersion int, module, handler string, request any, domain ...string) ([]byte, error) {
	var requestData []byte
	switch t := request.(type) {
	case string:
		requestData = convert.StringToBytes(t)
	case []byte:
		requestData = t
	default:
		v, err := json.Marshal(request)
		if err != nil {
			return nil, errors.Wrap(err, "marshal invoke request")
		}
		requestData = v
	}

	daprClient, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接
	daprAppId := namespace.Encapsulate(app)
	method := utils.GenerateMethod(apiVersion, module, handler, domain...)
	resp, err := daprClient.InvokeMethodWithContent(a.ctx, daprAppId, method, "post", &client.DataContent{
		ContentType: "application/json",
		Data:        requestData,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "dapr invoke method, appId: %s, method: %s", daprAppId, method)
	}

	return resp, nil
}
