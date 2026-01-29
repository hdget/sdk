package meta

import (
	"strconv"
	"sync"
)

const (
	KeyAppKey  = "hd-app-key"  // 应用ID
	KeySource  = "hd-source"   // 请求来源, 可以表示客户端，外部域，第三方API
	KeyRelease = "hd-release"  // 版本号
	KeyTid     = "hd-tid"      // tenant id
	KeyUid     = "hd-uid"      // user id
	KeyUsn     = "hd-usn"      // user sn
	KeyRoleIds = "hd-role-ids" // role ids
	KeyCaller  = "dapr-caller-app-id"
)

type MetaData interface {
	Set(key string, val any)
	Get(key string) (any, bool)
	AsGRPCMetaData() map[string][]string
	GetInt64(key string) int64
	GetString(key string) string
	GetInt64Slice(key string) []int64
}

type metaDataImpl struct {
	kvs map[string]any
	mu  sync.RWMutex
}

func New(kvs ...string) MetaData {
	c := &metaDataImpl{
		kvs: make(map[string]any),
		mu:  sync.RWMutex{},
	}

	if len(kvs) > 0 && len(kvs)%2 == 0 {
		c.mu.Lock()
		defer c.mu.Unlock()

		for i := 0; i < len(kvs); i += 2 {
			c.kvs[kvs[i]] = kvs[i+1]
		}
	}
	return c
}

func (impl *metaDataImpl) GetInt64Slice(key string) []int64 {
	if v, exists := impl.Get(key); exists {
		if val, ok := v.([]int64); ok {
			return val
		}
	}
	return nil
}

func (impl *metaDataImpl) GetInt64(key string) int64 {
	if v, exists := impl.Get(key); exists {
		if val, ok := v.(int64); ok {
			return val
		}
	}
	return 0
}

func (impl *metaDataImpl) GetString(key string) string {
	if v, exists := impl.Get(key); exists {
		if val, ok := v.(string); ok {
			return val
		}
	}
	return ""
}

func (impl *metaDataImpl) Get(key string) (any, bool) {
	impl.mu.RLock()
	defer impl.mu.RUnlock()
	val, exist := impl.kvs[key]
	return val, exist
}

func (impl *metaDataImpl) Set(key string, val any) {
	impl.mu.Lock()
	defer impl.mu.Unlock()
	impl.kvs[key] = val
}

func (impl *metaDataImpl) AsGRPCMetaData() map[string][]string {
	impl.mu.RLock()
	defer impl.mu.RUnlock()

	md := make(map[string][]string)
	for key, v := range impl.kvs {
		switch val := v.(type) {
		case string:
			md[key] = []string{val}
		case int64:
			md[key] = []string{strconv.FormatInt(val, 10)}
		case []int64:
			md[key] = make([]string, len(val))
			for i, vv := range val {
				md[key][i] = strconv.FormatInt(vv, 10)
			}
		default:
			// 对于未处理类型，使用fallback方案
			md[key] = []string{""}
		}
	}
	return md
}
