package dapr

import (
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
	"strings"
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
	appId := normalize(app)
	method := generateMethod(version, module, handler, clientName...)
	resp, err := daprClient.InvokeMethodWithContent(a.ctx, appId, method, "post", &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "dapr daprInvoke method, appId:%s, method: %s", appId, method)
	}

	return resp, nil
}

func generateMethod(version int, module, handler string, clientName ...string) string {
	tokens := []string{
		fmt.Sprintf("v%d", version),
		module,
		handler,
	}

	if len(clientName) > 0 && clientName[0] != "" {
		tokens = append(tokens, clientName[0])
	}

	return strings.ToLower(strings.Join(tokens, ":"))
}

//
//func parseClient(ctx context.Context) string {
//	md, ok := metadata.FromOutgoingContext(ctx)
//	if !ok {
//		return ""
//	}
//
//	values := md.Get(constant.MetaKeyClient)
//	if len(values) == 0 {
//		return ""
//	}
//	return values[0]
//}
