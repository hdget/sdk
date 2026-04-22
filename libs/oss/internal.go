package oss

import "time"

// internalConfigurer 配置接口，供 Option 函数使用
type internalConfigurer interface {
	setContentTypes(contentTypes []string)
	setMaxFileSize(size int64)
	setSignExpiresIn(duration time.Duration)
	setObjectACL(acl ObjectACL)
}

// WithContentTypes 设置允许的文件类型
func WithContentTypes(contentTypes []string) Option {
	return func(api API) {
		if configurer, ok := api.(internalConfigurer); ok && len(contentTypes) > 0 {
			configurer.setContentTypes(contentTypes)
		}
	}
}

// WithMaxFileSize 设置允许最大的文件大小
func WithMaxFileSize(size int64) Option {
	return func(api API) {
		if configurer, ok := api.(internalConfigurer); ok && size > 0 {
			configurer.setMaxFileSize(size)
		}
	}
}

// WithSignExpiresIn 设置签名过期时间
func WithSignExpiresIn(duration time.Duration) Option {
	return func(api API) {
		if configurer, ok := api.(internalConfigurer); ok && duration > 0 {
			configurer.setSignExpiresIn(duration)
		}
	}
}

// WithObjectACL 设置对象访问权限
func WithObjectACL(acl ObjectACL) Option {
	return func(api API) {
		if configurer, ok := api.(internalConfigurer); ok && acl != "" {
			configurer.setObjectACL(acl)
		}
	}
}
