package module

import (
	"github.com/hdget/utils"
	"github.com/hdget/utils/text"
)

var (
	truncateSize = 200
)

func truncate(data []byte) string {
	return text.Truncate(utils.BytesToString(data), truncateSize)
}
