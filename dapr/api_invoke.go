package dapr

import (
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/hdget/common/types"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
	"strings"
)

const ContentTypeJson = "application/json"

// Invoke 调用dapr服务
func (a apiImpl) Invoke(app string, moduleInfo *types.DaprModuleInfo, handler string, data any) ([]byte, error) {
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
	resp, err := daprClient.InvokeMethodWithContent(a.ctx, a.normalize(app), generateMethodName(moduleInfo, handler), "post", &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func generateMethodName(moduleInfo *types.DaprModuleInfo, handler string) string {
	if moduleInfo == nil {
		return ""
	}
	tokens := []string{
		fmt.Sprintf("v%d", moduleInfo.Version),
		moduleInfo.Name,
		handler,
	}
	// 去掉module后面的可能的module后缀
	if moduleInfo.Namespace != "" {
		tokens = append(tokens, moduleInfo.Namespace)
	}
	return strings.Join(tokens, ":")
}
