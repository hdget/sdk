package oss

import "time"

// InternalConfigurer 配置接口，供 Option 函数使用
type InternalConfigurer interface {
	SetContentTypes(contentTypes []string)
	SetMaxFileSize(size int64)
	SetSignExpiresIn(duration time.Duration)
	SetObjectACL(acl ObjectACL)
}

// WithContentTypes 设置允许的文件类型
func WithContentTypes(contentTypes []string) Option {
	return func(api API) {
		if configurer, ok := api.(InternalConfigurer); ok && len(contentTypes) > 0 {
			configurer.SetContentTypes(contentTypes)
		}
	}
}

// WithMaxFileSize 设置允许最大的文件大小
func WithMaxFileSize(size int64) Option {
	return func(api API) {
		if configurer, ok := api.(InternalConfigurer); ok && size > 0 {
			configurer.SetMaxFileSize(size)
		}
	}
}

// WithSignExpiresIn 设置签名过期时间
func WithSignExpiresIn(duration time.Duration) Option {
	return func(api API) {
		if configurer, ok := api.(InternalConfigurer); ok && duration > 0 {
			configurer.SetSignExpiresIn(duration)
		}
	}
}

// WithObjectACL 设置对象访问权限
func WithObjectACL(acl ObjectACL) Option {
	return func(api API) {
		if configurer, ok := api.(InternalConfigurer); ok && acl != "" {
			configurer.SetObjectACL(acl)
		}
	}
}
