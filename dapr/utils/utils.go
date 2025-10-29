package utils

import (
	"strconv"
	"strings"
)

func GenerateMethod(version int, module, handler string, apiEndpoint ...string) string {
	var builder strings.Builder

	// 预分配足够的内存以避免扩容开销
	// 估算基础长度：":v" (2) + 版本号（假设1-2位） + module长度 + handler长度 + 可能的AppEndpoint长度 + 分隔符":"
	estimatedLength := 2 + len(module) + len(handler) + len(":")
	for _, c := range apiEndpoint {
		estimatedLength += len(c)
	}
	builder.Grow(estimatedLength)

	// 直接写入，避免任何中间的字符串拼接
	builder.WriteString("v")
	builder.WriteString(strconv.Itoa(version))
	builder.WriteString(":")
	builder.WriteString(module)
	builder.WriteString(":")
	builder.WriteString(handler)

	if len(apiEndpoint) > 0 && apiEndpoint[0] != "" {
		builder.WriteString(":")
		builder.WriteString(apiEndpoint[0])
	}

	return strings.ToLower(builder.String())
}
