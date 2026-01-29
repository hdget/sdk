package namespace

import (
	"os"
	"strings"

	"github.com/hdget/sdk/common/constant"
)

const (
	separator = "-"
)

// Encapsulate 封装到命名空间
func Encapsulate(input string) string {
	if namespace, exists := os.LookupEnv(constant.EnvKeyNamespace); exists {
		var sb strings.Builder
		sb.Grow(len(namespace) + len(input) + 1)
		sb.WriteString(namespace)
		sb.WriteString(separator)
		sb.WriteString(input)
		return sb.String()
	}
	return input
}
