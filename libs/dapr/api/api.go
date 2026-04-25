package api

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
)

type DaprApi interface {
	Invoke(ctx context.Context, app string, apiVersion int, module string, handler string, data any, appCode ...string) ([]byte, error)
	Lock(ctx context.Context, lockStore, lockOwner, resource string, expiryInSeconds int) error
	Unlock(ctx context.Context, lockStore, lockOwner, resource string) error
	Publish(ctx context.Context, pubSubName, topic string, data interface{}, args ...bool) error
	SaveState(ctx context.Context, storeName, key string, value interface{}) error
	GetState(ctx context.Context, storeName, key string) ([]byte, error)
	DeleteState(ctx context.Context, storeName, key string) error
	GetConfigurationItems(ctx context.Context, configStore string, keys []string) (map[string]*client.ConfigurationItem, error)
	SubscribeConfigurationItems(ctx context.Context, configStore string, keys []string, handler client.ConfigurationHandleFunction) (string, error)
	GetBulkState(ctx context.Context, storeName string, keys any) (map[string][]byte, error)
}

type daprApiImpl struct {
}

var (
	newApi = sync.OnceValue[DaprApi](func() DaprApi {
		return &daprApiImpl{}
	})
)

func New() DaprApi {
	return newApi()
}

// InternalCall 内部调用, 不返回结果
func InternalCall(ctx context.Context, app string, version int, module, handler string, request ...any) error {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	_, err := New().Invoke(ctx, app, version, module, handler, req)
	if err != nil {
		return errors.Wrapf(err, "dapr internal call, app: %s, version: %d, module: %s, handler: %s, req: %v", app, version, module, handler, req)
	}

	return nil
}

// InternalInvoke 内部调用, 返回结果
func InternalInvoke[RESULT any](ctx context.Context, app string, version int, module, handler string, request ...any) (RESULT, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	var ret RESULT
	data, err := New().Invoke(ctx, app, version, module, handler, req)
	if err != nil {
		return ret, errors.Wrapf(err, "dapr internal invoke, app: %s, version: %d, module: %s, handler: %s, req: %v", app, version, module, handler, req)
	}

	err = json.Unmarshal(data, &ret)
	if err != nil {
		return ret, errors.Wrapf(err, "invalid dapr internal invoke result, app: %s, version: %d, module: %s, handler: %s, req: %v, ret: %v", app, version, module, handler, req, ret)
	}
	return ret, nil
}
