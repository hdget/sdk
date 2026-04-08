package validator

import "sync"

// 全局验证器实例（通过 Register 注册具体实现）
var (
	globalValidator Validator
	once            sync.Once
)

// Register 注册验证器实现（应用启动时调用一次）
func Register(v Validator) {
	once.Do(func() {
		globalValidator = v
	})
}

// Validate 使用全局验证器验证输入（便捷函数）
func Validate(input any) error {
	if globalValidator == nil {
		panic("validator not registered, call validator.Register() first")
	}
	return globalValidator.Validate(input)
}

// MustValidate 使用全局验证器验证，失败时 panic
func MustValidate(input any) {
	if globalValidator == nil {
		panic("validator not registered, call validator.Register() first")
	}
	globalValidator.MustValidate(input)
}
