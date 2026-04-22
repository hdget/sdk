package aliyun

import (
	"time"

	"github.com/hdget/sdk/libs/oss"
)

// 实现 oss 包 Option 函数所需的 setter 方法, 通过类型断言在 Option 函数中调用
// nolint:unused
func (impl *aliyunOssImpl) setContentTypes(contentTypes []string) {
	if len(contentTypes) > 0 {
		impl.allowContentTypes = contentTypes
	}
}

// nolint:unused
func (impl *aliyunOssImpl) setMaxFileSize(size int64) {
	if size > 0 {
		impl.maxFileSize = size
	}
}

// nolint:unused
func (impl *aliyunOssImpl) setSignExpiresIn(duration time.Duration) {
	if duration > 0 {
		impl.signExpiresIn = duration
	}
}

// nolint:unused
func (impl *aliyunOssImpl) setObjectACL(acl oss.ObjectACL) {
	if acl != "" {
		impl.objectACL = acl
	}
}
