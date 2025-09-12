package utils

import (
	"fmt"
	"github.com/hdget/common/constant"
	"os"
	"strings"
)

func Normalize(input string) string {
	if namespace, exists := os.LookupEnv(constant.EnvKeyNamespace); exists {
		return namespace + "_" + input
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
