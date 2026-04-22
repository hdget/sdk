package lib_ws

import (
	"slices"

	"github.com/gin-gonic/gin"
)

const keyMeta = "hd-meta"

// GetMetaKvs 获取meta的kv pairs
func GetMetaKvs(c *gin.Context) []string {
	return c.GetStringSlice(keyMeta)
}

// AddMetaKvs 添加信息到meta
func AddMetaKvs(c *gin.Context, kvs ...string) {
	if len(kvs)%2 == 1 {
		return
	}
	c.Set(keyMeta, slices.Concat(GetMetaKvs(c), kvs))
}
