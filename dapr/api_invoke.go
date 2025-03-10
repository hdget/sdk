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

// as we use gogoprotobuf which doesn't has protojson.RecvMessage interface
//var jsonpb = protojson.MarshalOptions{
//	EmitUnpopulated: true,
//}
//var jsonpbMarshaler = jsonpb.Marshaler{EmitDefaults: true}

// Invoke 调用dapr服务
func (a apiImpl) Invoke(app string, moduleVersion int, moduleName, handler string, data any) ([]byte, error) {
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

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer daprClient.Close()

	fullMethodName := getServiceInvocationName(moduleVersion, moduleName, handler)
	resp, err := daprClient.InvokeMethodWithContent(a.ctx, a.normalize(app), fullMethodName, "post", &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetServiceInvocationName 构造version:module:realMethod的方法名
func getServiceInvocationName(moduleVersion int, moduleName, handler string) string {
	return strings.Join([]string{fmt.Sprintf("v%d", moduleVersion), moduleName, handler}, ":")
}
