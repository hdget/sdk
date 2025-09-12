package biz

import (
	"github.com/hdget/common/protobuf"
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ErrCodeStart 业务逻辑错误代码开始值
const ErrCodeStart = 10000

// ErrCodeModuleRoot define error code module, e,g: 10000, 20000, 30000...
const (
	ErrCodeModuleRoot = ErrCodeStart * (1 + iota)
)

// define utils error code
const (
	_               = ErrCodeModuleRoot + iota // unknown error code
	ErrCodeInternal                            // internal error
)

type Error interface {
	error
	Code() int
}

type errorImpl struct {
	ErrCode int    `json:"code"`
	ErrMsg  string `json:"msg"`
}

// NewError new error with error code
func NewError[T constraints.Integer](code T, message string) Error {
	return &errorImpl{
		ErrCode: int(code),
		ErrMsg:  message,
	}
}

func ToGrpcError(err error) error {
	var be Error
	ok := errors.As(err, &be)
	if ok {
		st, _ := status.New(codes.Internal, "internal error").WithDetails(&protobuf.Error{
			Code: int32(be.Code()),
			Msg:  be.Error(),
		})
		return st.Err()
	}
	return err
}

// FromGrpcError 从grpc status error获取额外的错误信息
func FromGrpcError(err error) Error {
	if err == nil {
		return nil
	}

	cause := errors.Cause(err)
	st, ok := status.FromError(cause)
	if ok {
		details := st.Details()
		if len(details) > 0 {
			var pbErr protobuf.Error
			e := proto.Unmarshal(st.Proto().Details[0].GetValue(), &pbErr)
			if e == nil {
				return &errorImpl{
					ErrCode: int(pbErr.Code),
					ErrMsg:  pbErr.Msg,
				}
			}
		}
	}

	return &errorImpl{
		ErrCode: int(codes.Internal),
		ErrMsg:  err.Error(),
	}
}

func InternalError(message string) Error {
	return &errorImpl{
		ErrCode: ErrCodeInternal, // 10001为自定义的内部错误编码
		ErrMsg:  message,
	}
}

func (be errorImpl) Error() string {
	return be.ErrMsg
}

func (be errorImpl) Code() int {
	return be.ErrCode
}
