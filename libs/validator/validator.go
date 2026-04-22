package validator

import (
	"errors"
	"sync"
)

// 全局验证器实例（通过 Register 注册具体实现）
var (
	globalValidator Validator
	once            sync.Once
	errNotRegistered = errors.New("validator not registered, call validator.Register() first")
)

// Register 注册验证器实现（应用启动时调用一次）
func Register(v Validator) {
	once.Do(func() {
		globalValidator = v
	})
}

// Validate 使用全局验证器验证输入（便捷函数）
// 返回错误而非panic，让调用者决定如何处理
func Validate(input any) error {
	if globalValidator == nil {
		return errNotRegistered
	}
	return globalValidator.Validate(input)
}

// MustValidate 使用全局验证器验证，失败时 panic
// 注意: 此函数在验证器未注册或验证失败时会panic，仅在确定已注册的场景使用
func MustValidate(input any) {
	if globalValidator == nil {
		panic("validator not registered, call validator.Register() first")
	}
	globalValidator.MustValidate(input)
}
