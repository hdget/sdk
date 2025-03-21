package dapr

import (
	"fmt"
	"github.com/hdget/utils/convert"
	"github.com/hdget/utils/text"
	"strings"
)

var (
	truncateSize = 200
)

func truncate(data []byte) string {
	return text.Truncate(convert.BytesToString(data), truncateSize)
}

// buildServiceInvocationName 构造version:module:realMethod的方法名
func buildServiceInvocationName(moduleVersion int, moduleName, handler string) string {
	return strings.Join([]string{fmt.Sprintf("v%d", moduleVersion), moduleName, handler}, ":")
}
