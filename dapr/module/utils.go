package module

import (
	"github.com/hdget/utils/convert"
	"github.com/hdget/utils/text"
)

var (
	truncateSize = 200
)

func truncate(data []byte) string {
	return text.Truncate(convert.BytesToString(data), truncateSize)
}
