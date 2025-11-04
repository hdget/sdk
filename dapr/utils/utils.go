package utils

import (
	"strconv"
	"strings"
)

func GenerateMethod(apiVersion int, module, handler string, client ...string) string {
	var builder strings.Builder

	// 预分配足够的内存以避免扩容开销
	// e,g:
	// http调用：v1:client:module:handler
	// 内部直接调用：v1:module:handler
	// 估算基础长度：v (1) + 版本号（假设1位）+ 2分隔符":" + module长度 + handler长度
	estimatedLength := 4 + len(module) + len(handler)
	if len(client) > 0 {
		estimatedLength += len(client[0]) + 1
	}
	builder.Grow(estimatedLength)

	// 直接写入，避免任何中间的字符串拼接
	builder.WriteString("v")
	builder.WriteString(strconv.Itoa(apiVersion))
	if len(client) > 0 {
		builder.WriteString(":")
		builder.WriteString(client[0])
	}
	builder.WriteString(":")
	builder.WriteString(module)
	builder.WriteString(":")
	builder.WriteString(handler)

	return strings.ToLower(builder.String())
}
