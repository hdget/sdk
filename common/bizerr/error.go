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
	defaultErrCode   = 10001
	defaultErrReason = "INTERNAL_ERROR"
)

// New error with error code
func New[T constraints.Integer](code T, reason, message string, kvs ...any) Error {
	return &bizErrorImpl{
		ErrCode:   int(code),
		ErrMsg:    message,
		ErrReason: reason,
		ErrDetail: parseKvs(kvs...),
	}
}

func FromProto(enum protoreflect.Enum, message string, kvs ...any) Error {
	if enum == nil {
		return InternalError(message, kvs...)
	}

	reason := defaultErrReason
	if ed := enum.Descriptor().Values().ByNumber(enum.Number()); ed != nil {
		reason = string(ed.Name())
	}

	return &bizErrorImpl{
		ErrCode:   int(enum.Number()),
		ErrReason: reason,
		ErrMsg:    message,
		ErrDetail: parseKvs(kvs...),
	}
}

func InternalError(message string, kvs ...any) Error {
	return New(defaultErrCode, defaultErrReason, message, kvs...)
}

func ToGrpcError(err error) error {
	var be *bizErrorImpl
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

	st, _ := status.New(codes.Unknown, be.Error()).WithDetails(pbErr)
	return st.Err()
}

// FromGrpcError 从grpc status error获取额外的错误信息
func FromGrpcError(err error) Error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(errors.Cause(err))
	if !ok {
		return &bizErrorImpl{
			ErrCode:   defaultErrCode,
			ErrReason: defaultErrReason,
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
		ErrReason: defaultErrReason,
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

	parsedMap := parseKvs(kvs...)
	if len(parsedMap) > 0 {
		cp.ErrDetail = make(map[string]any, len(parsedMap))
		for k, v := range parsedMap {
			cp.ErrDetail[k] = v
		}
	} else {
		cp.ErrDetail = nil
	}

	return &cp
}

func (be *bizErrorImpl) Detail() map[string]any {
	return be.ErrDetail
}

func parseKvs(kvs ...any) map[string]any {
	if len(kvs) == 0 {
		return nil
	}

	// 如果是 map / struct，直接走结构化转换
	if len(kvs) == 1 {
		if m, ok := kvs[0].(map[string]any); ok {
			return m
		}
	}

	// kvs key-value pair 模式
	if len(kvs)%2 != 0 {
		// 可以选择 panic / ignore / log
		return map[string]any{"_err": "invalid detail"}
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
