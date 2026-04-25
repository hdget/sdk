package rabbitmq

import (
	"github.com/hdget/sdk/common/provider"
)

// rabbitmqProvider
// Note: most codes comes from https://github.com/ThreeDotsLabs/watermill-amqp
type rabbitmqProvider struct {
	config *RabbitMqConfig
	logger provider.Logger
}

func New(configProvider provider.Config, logger provider.Logger) (provider.MessageQueue, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	return &rabbitmqProvider{config: config, logger: logger}, nil
}

func (r rabbitmqProvider) Init(args ...any) error {
	// RabbitMQ provider不需要额外的初始化步骤
	// 配置已在New()中完成验证
	return nil
}

func (r rabbitmqProvider) NewPublisher(name string, args ...*provider.PublisherOption) (provider.MessageQueuePublisher, error) {
	option := provider.DefaultPublisherOption
	if len(args) > 0 {
		option = args[0]
	}

	publisherOptions := make([]publisherOption, 0)
	if option.PublishDelayMessage {
		publisherOptions = append(publisherOptions, withPublisherDelayTopology())
	}

	return newPublisher(name, r.config, r.logger, publisherOptions...)
}

func (r rabbitmqProvider) NewSubscriber(name string, args ...*provider.SubscriberOption) (provider.MessageQueueSubscriber, error) {
	option := provider.DefaultSubscriberOption
	if len(args) > 0 {
		option = args[0]
	}

	subscriberOptions := make([]subscriberOption, 0)
	if option.SubscribeDelayMessage {
		subscriberOptions = append(subscriberOptions, withSubscriberDelayTopology())
	}

	return newSubscriber(name, r.config, r.logger, subscriberOptions...)
}

func (r rabbitmqProvider) GetCapability() provider.Capability {
	return Capability
}
