package module

import (
	"context"
	"encoding/json"

	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/common/biz"
	"github.com/hdget/common/types"
	"github.com/hdget/sdk/dapr/api"
	"github.com/hdget/sdk/dapr/utils"
	"github.com/hdget/utils/convert"
	panicUtils "github.com/hdget/utils/panic"
	reflectUtils "github.com/hdget/utils/reflect"
)

type invocationHandler interface {
	GetAlias() string
	GetName() string
	GetInvokeName() string                                                         // 调用名字
	GetInvokeFunction(logger types.LoggerProvider) common.ServiceInvocationHandler // 具体的调用函数
}

type invocationHandlerImpl struct {
	module Module
	// handler的别名，
	// 如果DiscoverHandlers调用, 会将函数名作为入参，matchFunction的返回值当作别名，缺省是去除Handler后缀并小写
	// 如果RegisterHandlers调用，会直接用map的key值当为别名
	handlerAlias string
	handlerName  string             // 调用函数名
	fn           InvocationFunction // 调用函数
}

type InvocationFunction func(ctx biz.Context, data []byte) (any, error)
type HandlerMatcher func(methodName string) (string, bool) // 传入receiver.methodName, 判断是否匹配，然后取出处理后的handlerName

func (h invocationHandlerImpl) GetAlias() string {
	return h.handlerAlias
}

func (h invocationHandlerImpl) GetName() string {
	return h.handlerName
}

func (h invocationHandlerImpl) GetInvokeName() string {
	mInfo := h.module.GetInfo()
	return utils.GenerateMethod(mInfo.ApiVersion, mInfo.Name, h.handlerAlias, mInfo.Dir)
}

func (h invocationHandlerImpl) GetInvokeFunction(logger types.LoggerProvider) common.ServiceInvocationHandler {
	return func(ctx context.Context, event *common.InvocationEvent) (*common.Content, error) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				panicUtils.RecordErrorStack(h.module.GetApp())
			}
		}()

		result, err := h.fn(biz.NewFromIncomingGrpcContext(ctx), event.Data)
		if err != nil {
			mInfo := h.module.GetInfo()
			logger.Error("service invoke", "domain", mInfo.Dir, "module", mInfo.Name, "Handler", reflectUtils.GetFuncName(h.fn), "err", err, "req", truncate(event.Data))
			return h.replyError(err)
		}

		return h.replySuccess(event, result)
	}
}

func (h invocationHandlerImpl) replyError(err error) (*common.Content, error) {
	return nil, biz.ToGrpcError(err)
}

func (h invocationHandlerImpl) replySuccess(event *common.InvocationEvent, result any) (*common.Content, error) {
	var err error
	var data []byte
	switch t := result.(type) {
	case string:
		data = convert.StringToBytes(t)
	case []byte:
		data = t
	default:
		data, err = json.Marshal(result)
		if err != nil {
			return nil, err
		}
	}

	return &common.Content{
		ContentType: api.ContentTypeJson,
		Data:        data,
		DataTypeURL: event.DataTypeURL,
	}, nil
}
