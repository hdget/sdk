package api

import (
	daprclient "github.com/dapr/go-sdk/client"
	"github.com/hdget/sdk/common/namespace"
	"github.com/hdget/sdk/libs/dapr/localutils"
	"github.com/hdget/utils"
	"github.com/pkg/errors"
)

const ContentTypeJson = "application/json"

// Invoke 调用dapr服务
func (a apiImpl) Invoke(app string, apiVersion int, module, handler string, request any, client ...string) ([]byte, error) {
	requestData, err := utils.ToBytes(request)
	if err != nil {
		return nil, errors.Wrap(err, "marshal invoke request")
	}

	c, err := daprclient.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接
	daprAppId := namespace.Encapsulate(app)
	method := localutils.GenerateMethod(apiVersion, module, handler, client...)
	resp, err := c.InvokeMethodWithContent(a.ctx, daprAppId, method, "post", &daprclient.DataContent{
		ContentType: "application/json",
		Data:        requestData,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "dapr invoke method, appId: %s, method: %s", daprAppId, method)
	}

	return resp, nil
}
