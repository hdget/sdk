package api

import (
	"context"
	"encoding/json"

	"github.com/dapr/go-sdk/client"
	"github.com/hdget/common/biz"
	"github.com/pkg/errors"
)

type APIer interface {
	Invoke(app string, apiVersion int, module string, handler string, data any, source ...string) ([]byte, error)
	Lock(lockStore, lockOwner, resource string, expiryInSeconds int) error
	Unlock(lockStore, lockOwner, resource string) error
	Publish(pubSubName, topic string, data interface{}, args ...bool) error
	SaveState(storeName, key string, value interface{}) error
	GetState(storeName, key string) ([]byte, error)
	DeleteState(storeName, key string) error
	GetConfigurationItems(configStore string, keys []string) (map[string]*client.ConfigurationItem, error)
	SubscribeConfigurationItems(ctx context.Context, configStore string, keys []string, handler client.ConfigurationHandleFunction) (string, error)
	GetBulkState(storeName string, keys any) (map[string][]byte, error)
}

type apiImpl struct {
	ctx context.Context
}

func New(ctx biz.Context) APIer {
	return &apiImpl{
		ctx: biz.NewOutgoingGrpcContext(ctx),
	}
}

// InternalCall 内部调用, 不返回结果
func InternalCall(ctx biz.Context, app string, version int, module, handler string, request ...any) error {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	_, err := New(ctx).Invoke(app, version, module, handler, req)
	if err != nil {
		return errors.Wrapf(err, "dapr internal call, app: %s, version: %d, module: %s, handler: %s, req: %v", app, version, module, handler, req)
	}

	return nil
}

// InternalInvoke 内部调用, 返回结果
func InternalInvoke[RESULT any](ctx biz.Context, app string, version int, module, handler string, request ...any) (*RESULT, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	data, err := New(ctx).Invoke(app, version, module, handler, req)
	if err != nil {
		return nil, errors.Wrapf(err, "dapr internal invoke, app: %s, version: %d, module: %s, handler: %s, req: %v", app, version, module, handler, req)
	}

	var ret RESULT
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid dapr internal invoke result, app: %s, version: %d, module: %s, handler: %s, req: %v, ret: %v", app, version, module, handler, req, ret)
	}
	return &ret, nil
}
