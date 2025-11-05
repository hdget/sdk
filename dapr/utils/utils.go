package utils

import (
	"strconv"
	"strings"
)

func GenerateMethod(apiVersion int, module, handler string, domain ...string) string {
	var builder strings.Builder

	var accessDomain string
	if len(domain) > 0 && domain[0] != "" {
		accessDomain = domain[0]
	}

	// 预分配足够的内存以避免扩容开销
	// e,g:
	// http调用：v1:domain:module:handler
	// 内部直接调用：v1:module:handler
	// 估算基础长度：v (1) + 版本号（假设1位）+ 2分隔符":" + module长度 + handler长度
	estimatedLength := 4 + len(module) + len(handler)
	if accessDomain != "" {
		estimatedLength += len(accessDomain) + 1
	}
	builder.Grow(estimatedLength)

	// 直接写入，避免任何中间的字符串拼接
	builder.WriteString("v")
	builder.WriteString(strconv.Itoa(apiVersion))
	if accessDomain != "" {
		builder.WriteString(":")
		builder.WriteString(accessDomain)
	}
	builder.WriteString(":")
	builder.WriteString(module)
	builder.WriteString(":")
	builder.WriteString(handler)

	return strings.ToLower(builder.String())
}
