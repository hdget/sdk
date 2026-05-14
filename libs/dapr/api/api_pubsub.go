package api

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/sdk/common/namespace"
	"github.com/pkg/errors"
)

type event struct {
	Subscription *common.Subscription
	Handler      common.TopicEventHandler
}

const (
	defaultPubSub       = "pubsub"
	sysEventTopicPrefix = "sys:event"
)

func NewEvent(pubsubName, topic string, handler common.TopicEventHandler, args ...bool) event {
	metaOptions := getPublishMetaOptions(args...)
	return event{
		Subscription: &common.Subscription{
			PubsubName: pubsubName,
			Topic:      topic,
			Metadata:   metaOptions,
		},
		Handler: handler,
	}
}

// Publish 发布消息
// isRawPayLoad 发送原始的消息，非cloudevent message
func (a daprApiImpl) Publish(ctx context.Context, pubSubName, topic string, data interface{}, args ...bool) error {
	c, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer c.Close()

	var opt client.PublishEventOption
	metaOptions := getPublishMetaOptions(args...)
	if metaOptions != nil {
		opt = client.PublishEventWithMetadata(metaOptions)
		err = c.PublishEvent(ctx, namespace.Encapsulate(pubSubName), topic, data, opt)
	} else {
		err = c.PublishEvent(ctx, namespace.Encapsulate(pubSubName), topic, data)
	}

	if err != nil {
		return err
	}

	return nil
}

// PublishSysEvent 发布系统事件
func PublishSysEvent[EventKind fmt.Stringer](ctx context.Context, kind EventKind, data any, pubSubName ...string) error {
	pubsub := defaultPubSub
	if len(pubSubName) > 0 {
		pubsub = pubSubName[0]
	}
	return New().Publish(ctx, pubsub, GetSysEventTopic(kind), data)
}

func GetSysEventTopic[EventKind fmt.Stringer](kind EventKind) string {
	return fmt.Sprintf("%s:%s", sysEventTopicPrefix, kind.String())
}

func getPublishMetaOptions(args ...bool) map[string]string {
	isRawPayLoad := false
	if len(args) > 0 {
		isRawPayLoad = args[0]
	}

	if isRawPayLoad {
		return map[string]string{"rawPayload": "true"}
	}
	return nil
}
