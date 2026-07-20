package aliyun

import (
	"time"

	"github.com/hdget/sdk/libs/oss"
)

// 实现 oss.InternalConfigurer 接口的方法
// nolint:unused
func (impl *aliyunOssImpl) SetContentTypes(contentTypes []string) {
	if len(contentTypes) > 0 {
		impl.allowContentTypes = contentTypes
	}
}

// nolint:unused
func (impl *aliyunOssImpl) SetMaxFileSize(size int64) {
	if size > 0 {
		impl.maxFileSize = size
	}
}

// nolint:unused
func (impl *aliyunOssImpl) SetSignExpiresIn(duration time.Duration) {
	if duration > 0 {
		impl.signExpiresIn = duration
	}
}

// nolint:unused
func (impl *aliyunOssImpl) SetObjectACL(acl oss.ObjectACL) {
	if acl != "" {
		impl.objectACL = acl
	}
}
