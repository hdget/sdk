package dapr

import (
	"context"
	"github.com/cenkalti/backoff/v4"
	"github.com/hdget/common/intf"
	"github.com/hdget/common/types"
	panicUtils "github.com/hdget/utils/panic"
	"time"
)

type DelayEventHandler interface {
	GetTopic() string
	Handle(ctx context.Context, logger intf.LoggerProvider, msgChan <-chan *types.Message)
}

type delayEventHandlerImpl struct {
	module DelayEventModule
	topic  string             // 订阅主题
	fn     DelayEventFunction // 调用函数
}

type DelayEventFunction func(message []byte) (retry bool, err error)

func (h delayEventHandlerImpl) GetTopic() string {
	// 如果使用的rabbitmq, 则第一个为实际topic, 第二个值为exchange
	return h.topic
}

// Handle
// err: nil 只要错误为空，则消息成功消费, 不管retry的值为什么样
// err: not nil + retry: false 打印DROP status消息
// err: not nil + retry: true  进行重试，最后重试次数结束, 打印日志
func (h delayEventHandlerImpl) Handle(ctx context.Context, logger intf.LoggerProvider, msgChan <-chan *types.Message) {
	// 挂载defer函数
	defer func() {
		if r := recover(); r != nil {
			panicUtils.RecordErrorStack(h.module.GetApp())
		}
	}()

LOOP:
	for {
		select {
		case <-ctx.Done():
			logger.Debug("shutdown delay event handler", "topic", h.GetTopic())
			break LOOP
		case msg := <-msgChan:
			retry, err := h.fn(msg.Payload)
			if err == nil {
				msg.Ack()
			} else {
				if !retry { // err != nil && retry == false
					logger.Error("drop delay event", "err", err, "data", truncate(msg.Payload))
					msg.Ack()
				} else { // err != nil && retry == true
					nextBackOff := h.module.GetBackOffPolicy().NextBackOff()
					if nextBackOff == backoff.Stop {
						logger.Error("drop delay event after retried many times", "err", err, "data", truncate(msg.Payload))
						msg.Ack()
						h.module.GetBackOffPolicy().Reset()
					} else {
						time.Sleep(nextBackOff)
						logger.Error("retry delay event", "err", err, "data", truncate(msg.Payload))
						msg.Nack()
					}
				}
			}
		}
	}
}
