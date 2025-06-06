package dapr

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/common/intf"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/sdk/bizerr"
	"github.com/hdget/utils/convert"
	panicUtils "github.com/hdget/utils/panic"
	reflectUtils "github.com/hdget/utils/reflect"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type invocationHandler interface {
	GetAlias() string
	GetName() string
	GetInvokeName() string                                                        // 调用名字
	GetInvokeFunction(logger intf.LoggerProvider) common.ServiceInvocationHandler // 具体的调用函数
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

type InvocationFunction func(ctx context.Context, event *common.InvocationEvent) (any, error)
type HandlerMatcher func(methodName string) (string, bool) // 传入receiver.methodName, 判断是否匹配，然后取出处理后的handlerName

func (h invocationHandlerImpl) GetAlias() string {
	return h.handlerAlias
}

func (h invocationHandlerImpl) GetName() string {
	return h.handlerName
}

func (h invocationHandlerImpl) GetInvokeName() string {
	mInfo := h.module.GetModuleInfo()
	return generateMethod(mInfo.Version, mInfo.Name, h.handlerAlias, mInfo.Client)
}

func (h invocationHandlerImpl) GetInvokeFunction(logger intf.LoggerProvider) common.ServiceInvocationHandler {
	return func(ctx context.Context, event *common.InvocationEvent) (*common.Content, error) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				panicUtils.RecordErrorStack(h.module.GetApp())
			}
		}()

		result, err := h.fn(ctx, event)
		if err != nil {
			mInfo := h.module.GetModuleInfo()
			logger.Error("service daprInvoke", "client", mInfo.Client, "module", mInfo.Name, "Handler", reflectUtils.GetFuncName(h.fn), "err", err, "req", truncate(event.Data))
			return h.replyError(err)
		}

		return h.replySuccess(event, result)
	}
}

func (h invocationHandlerImpl) replyError(err error) (*common.Content, error) {
	var be *bizerr.BizError
	ok := errors.As(err, &be)
	if ok {
		st, _ := status.New(codes.Internal, "internal error").WithDetails(&protobuf.Error{
			Code: be.Code,
			Msg:  be.Msg,
		})
		return nil, st.Err()
	}
	return nil, err
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
		ContentType: ContentTypeJson,
		Data:        data,
		DataTypeURL: event.DataTypeURL,
	}, nil
}
