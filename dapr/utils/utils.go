package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/hdget/common/constant"
)

const (
	// 命名空间和真实名字的分割符号
	namespaceNameSeparator = "-"
)

func Normalize(input string) string {
	if namespace, exists := os.LookupEnv(constant.EnvKeyNamespace); exists {
		var sb strings.Builder
		sb.Grow(len(namespace) + len(input) + 1)
		sb.WriteString(namespace)
		sb.WriteString(namespaceNameSeparator)
		sb.WriteString(input)
		return sb.String()
	}
	return input
}

func GenerateMethod(version int, module, handler string, client ...string) string {
	tokens := []string{
		fmt.Sprintf("v%d", version),
		module,
		handler,
	}

	if len(client) > 0 && client[0] != "" {
		tokens = append(tokens, client[0])
	}

	return strings.ToLower(strings.Join(tokens, ":"))
}
