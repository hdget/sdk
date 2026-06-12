package bizerr

import (
	"github.com/hdget/sdk/common/protobuf"
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

type Error interface {
	error
	Code() int
	Reason() string
	Detail() map[string]any
	WithDetail(kvs ...any) Error
}

type bizErrorImpl struct {
	ErrCode   int            `json:"code"`
	ErrReason string         `json:"reason"`
	ErrMsg    string         `json:"msg"`
	ErrDetail map[string]any `json:"detail,omitempty"`
}

const (
	internalCode   = 10001
	internalReason = "INTERNAL_ERROR"
)

// New error with error code
func New[T constraints.Integer](code T, reason, message string, kvs ...any) Error {
	return &bizErrorImpl{
		ErrCode:   int(code),
		ErrMsg:    message,
		ErrReason: reason,
		ErrDetail: cloneMap(parseKvs(kvs...)),
	}
}

func FromProto(enum protoreflect.Enum, message string, kvs ...any) Error {
	if enum == nil {
		return InternalError(message, kvs...)
	}

	reason := internalReason
	// 增加中间值的 nil 检查
	if desc := enum.Descriptor(); desc != nil {
		if values := desc.Values(); values != nil {
			if v := values.ByNumber(enum.Number()); v != nil {
				reason = string(v.Name())
			}
		}
	}

	return &bizErrorImpl{
		ErrCode:   int(enum.Number()),
		ErrReason: reason,
		ErrMsg:    message,
		ErrDetail: cloneMap(parseKvs(kvs...)),
	}
}

func InternalError(message string, kvs ...any) Error {
	return New(internalCode, internalReason, message, kvs...)
}

func ToGrpcError(err error) error {
	var be Error
	if !errors.As(err, &be) {
		return err
	}

	pbErr := &protobuf.Error{
		Code:   int32(be.Code()),
		Msg:    be.Error(),
		Reason: be.Reason(),
	}

	if detail := be.Detail(); len(detail) > 0 {
		if pbDetail, e := structpb.NewStruct(detail); e == nil {
			pbErr.Detail = pbDetail
		}
	}

	st, _ := status.New(codes.Code(be.Code()), be.Error()).WithDetails(pbErr)
	return st.Err()
}

// FromGrpcError 从grpc status error获取额外的错误信息
func FromGrpcError(err error) Error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return &bizErrorImpl{
			ErrCode:   internalCode,
			ErrReason: internalReason,
			ErrMsg:    err.Error(),
		}
	}

	for _, d := range st.Details() {
		if pbErr, ok := d.(*protobuf.Error); ok {
			var detail map[string]any
			if pbErr.Detail != nil {
				detail = pbErr.Detail.AsMap()
			}

			return &bizErrorImpl{
				ErrCode:   int(pbErr.Code),
				ErrReason: pbErr.Reason,
				ErrMsg:    pbErr.Msg,
				ErrDetail: detail,
			}
		}
	}

	return &bizErrorImpl{
		ErrCode:   int(st.Code()),
		ErrReason: internalReason,
		ErrMsg:    st.Message(),
	}
}

func (be *bizErrorImpl) Error() string {
	return be.ErrMsg
}

func (be *bizErrorImpl) Code() int {
	return be.ErrCode
}

func (be *bizErrorImpl) Reason() string {
	return be.ErrReason
}

func (be *bizErrorImpl) WithDetail(kvs ...any) Error {
	cp := *be

	detail := cloneMap(be.ErrDetail)
	for k, v := range parseKvs(kvs...) {
		detail[k] = v
	}

	cp.ErrDetail = detail
	return &cp
}

func (be *bizErrorImpl) Detail() map[string]any {
	return cloneMap(be.ErrDetail)
}

func parseKvs(kvs ...any) map[string]any {
	if len(kvs) == 0 {
		return nil
	}

	// 如果是 map / struct，直接走结构化转换
	if len(kvs) == 1 {
		if m, ok := kvs[0].(map[string]any); ok {
			return cloneMap(m)
		}
	}

	// kvs key-value pair 模式
	if len(kvs)%2 != 0 {
		// 可以选择 panic / ignore / log
		return nil
	}

	out := make(map[string]any, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			// key 不是 string，跳过或报错
			continue
		}
		out[key] = kvs[i+1]
	}
	return out
}

func cloneMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}

	dst := make(map[string]any, len(src))

	for k, v := range src {
		dst[k] = v
	}

	return dst
}
