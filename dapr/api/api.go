package api

import (
	"context"
	"encoding/json"

	"github.com/dapr/go-sdk/client"
	"github.com/hdget/common/biz"
)

type APIer interface {
	Invoke(app string, version int, module string, handler string, data any, client ...string) ([]byte, error)
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

func Call[RESULT any](ctx biz.Context, app string, version int, module, handler string, request ...any) (*RESULT, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	data, err := New(ctx).Invoke(app, version, module, handler, req)
	if err != nil {
		return nil, err
	}

	var ret RESULT
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
