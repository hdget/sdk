package utils

import (
	"strconv"
	"strings"
)

func GenerateMethod(apiVersion int, module, handler string, source ...string) string {
	var builder strings.Builder

	var origin string
	if len(source) > 0 && source[0] != "" {
		origin = source[0]
	}

	// 预分配足够的内存以避免扩容开销
	// e,g:
	// 外部调用：source:v1:module:handler
	// 内部调用：v1:module:handler
	// 估算基础长度：v (1) + 版本号（假设1位）+ 2分隔符":" + module长度 + handler长度
	estimatedLength := 4 + len(module) + len(handler)
	if origin != "" {
		estimatedLength += len(origin) + 1
	}
	builder.Grow(estimatedLength)

	if origin != "" {
		builder.WriteString(origin)
		builder.WriteString(":")
	}

	// 直接写入，避免任何中间的字符串拼接
	builder.WriteString("v")
	builder.WriteString(strconv.Itoa(apiVersion))
	builder.WriteString(":")
	builder.WriteString(module)
	builder.WriteString(":")
	builder.WriteString(handler)

	return strings.ToLower(builder.String())
}
