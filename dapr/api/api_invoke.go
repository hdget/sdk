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
func (a apiImpl) Invoke(app string, version int, module, handler string, data any, clientName ...string) ([]byte, error) {
	var value []byte
	switch t := data.(type) {
	case string:
		value = convert.StringToBytes(t)
	case []byte:
		value = t
	default:
		v, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "marshal daprInvoke data")
		}
		value = v
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
	method := utils.GenerateMethod(version, module, handler, clientName...)
	resp, err := daprClient.InvokeMethodWithContent(a.ctx, daprAppId, method, "post", &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "dapr daprInvoke method, daprAppId:%s, method: %s", daprAppId, method)
	}

	return resp, nil
}
