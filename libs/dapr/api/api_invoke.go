package api

import (
	"context"

	"github.com/dapr/go-sdk/client"
	"github.com/hdget/sdk/common/bizctx"
	"github.com/hdget/sdk/common/namespace"
	"github.com/hdget/sdk/libs/dapr/localutils"
	"github.com/hdget/utils"
	"github.com/pkg/errors"
)

const ContentTypeJson = "application/json"

// Invoke 调用dapr服务
func (a daprApiImpl) Invoke(ctx context.Context, app string, apiVersion int, module, handler string, request any, appCode ...string) ([]byte, error) {
	requestData, err := utils.ToBytes(request)
	if err != nil {
		return nil, errors.Wrap(err, "marshal invoke request")
	}

	c, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr appCode")
	}
	if c == nil {
		return nil, errors.New("dapr appCode is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接
	daprAppId := namespace.Encapsulate(app)
	method := localutils.GenerateMethod(apiVersion, module, handler, appCode...)
	resp, err := c.InvokeMethodWithContent(bizctx.NewOutgoingGrpcContext(ctx), daprAppId, method, "post", &client.DataContent{
		ContentType: "application/json",
		Data:        requestData,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "dapr invoke method, appId: %s, method: %s", daprAppId, method)
	}

	return resp, nil
}
