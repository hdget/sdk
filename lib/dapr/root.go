package dapr

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
)

// InvokeOnce 调用一次
func InvokeOnce(appId, methodName string, data interface{}) ([]byte, error) {
	var value []byte
	switch t := data.(type) {
	case string:
		value = utils.StringToBytes(t)
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
	defer daprClient.Close()

	content := &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	}

	resp, err := daprClient.InvokeMethodWithContent(context.Background(), appId, methodName, "post", content)
	if err != nil {
		return nil, errors.Wrapf(err, "dapr invoke method with content, app:%s, method: %s, content: %s", appId, methodName, utils.BytesToString(value))
	}

	return resp, nil
}

// Invoke 需要传入daprClient去调用
func Invoke(daprClient client.Client, appId, methodName string, data interface{}) ([]byte, error) {
	var value []byte
	switch t := data.(type) {
	case string:
		value = utils.StringToBytes(t)
	default:
		v, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "marshal invoke data")
		}
		value = v
	}

	content := &client.DataContent{
		ContentType: "application/json",
		Data:        value,
	}

	ret, err := daprClient.InvokeMethodWithContent(context.Background(), appId, methodName, "post", content)
	if err != nil {
		return nil, errors.Wrapf(err, "dapr invoke method with content, app:%s, method: %s, content: %s", appId, methodName, utils.BytesToString(value))
	}

	return ret, nil
}