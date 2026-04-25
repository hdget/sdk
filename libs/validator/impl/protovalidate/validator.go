// Package validator 提供基于 protovalidate 的声明式参数验证功能
package protovalidate

import (
	"bytes"

	"buf.build/go/protovalidate"
	"github.com/hdget/sdk/libs/validator"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type protoValidatorImpl struct {
	protovalidate.Validator
}

// New 获取验证器实例
func New() validator.Validator {
	v, err := protovalidate.New()
	if err != nil {
		panic(err)
	}

	return &protoValidatorImpl{
		Validator: v,
	}
}

// Validate 验证 protobuf 消息
// 返回验证错误，如果验证通过则返回 nil
func (impl protoValidatorImpl) Validate(input any) error {
	if v, ok := input.(proto.Message); ok {
		return wrapError(impl.Validator.Validate(v))
	}
	return errors.New("input is not a proto message")
}

// MustValidate 验证 protobuf 消息，如果验证失败则 panic
func (impl protoValidatorImpl) MustValidate(input any) {
	if v, ok := input.(proto.Message); ok {
		if err := impl.Validator.Validate(v); err != nil {
			panic(err)
		}
	}
}

// WrapError 将 protovalidate 错误转换为 ValidationError
func wrapError(err error) error {
	if err == nil {
		return nil
	}

	var valErr *protovalidate.ValidationError
	if errors.As(err, &valErr) && len(valErr.Violations) > 0 {
		var buf bytes.Buffer
		for _, violation := range valErr.Violations {
			buf.WriteString(violation.String())
			buf.WriteString(",")
		}
		return errors.New(buf.String())
	}
	return err
}
