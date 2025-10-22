package utils

import (
	"fmt"
	"strings"
)

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
