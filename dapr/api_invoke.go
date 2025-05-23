package dapr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/hdget/common/constant"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"strings"
)

const ContentTypeJson = "application/json"

// Invoke 内部调用
func (a apiImpl) Invoke(app string, version int, module, handler string, data any) ([]byte, error) {
	return a.daprInvoke(app, version, module, handler, data, true)
}

// ExternalInvoke 外部调用
func (a apiImpl) ExternalInvoke(app string, version int, module, handler string, data any) ([]byte, error) {
	return a.daprInvoke(app, version, module, handler, data, false)
}

// daprInvoke 调用dapr服务
func (a apiImpl) daprInvoke(app string, version int, module, handler string, data any, ignoreClient bool) ([]byte, error) {
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
	var appId, method string
	appId = normalize(app)
	if ignoreClient {
		method = generateMethod(version, module, handler, "")
	} else {
		method = generateMethod(version, module, handler, getClient(a.ctx))
	}

	resp, err := daprClient.InvokeMethodWithContent(a.ctx, appId, method, "post", &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "dapr daprInvoke method, appId:%s, method: %s", appId, method)
	}

	return resp, nil
}

func getClient(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(constant.MetaKeyClient)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func generateMethod(version int, module, handler, client string) string {
	tokens := []string{
		fmt.Sprintf("v%d", version),
		module,
		handler,
	}

	if client != "" {
		tokens = append(tokens, client)
	}

	return strings.ToLower(strings.Join(tokens, ":"))
}
