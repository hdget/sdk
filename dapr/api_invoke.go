package dapr

import (
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"strings"
)

const ContentTypeJson = "application/json"

// Invoke 调用dapr服务
func (a apiImpl) Invoke(app string, version int, module, handler string, data any) ([]byte, error) {
	var value []byte
	switch t := data.(type) {
	case string:
		value = convert.StringToBytes(t)
	case []byte:
		value = t
	default:
		v, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "marshal invoke data")
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
	method := generateMethodName(version, module, handler, cast.ToString(a.ctx.Value(MetaKeyAppId)))
	resp, err := daprClient.InvokeMethodWithContent(a.ctx, a.normalize(app), method, "post", &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func generateMethodName(version int, module, handler, appId string) string {
	tokens := []string{
		fmt.Sprintf("v%d", version),
		module,
		handler,
	}

	if appId != "" {
		tokens = append(tokens, appId)
	}

	return strings.ToLower(strings.Join(tokens, ":"))
}
