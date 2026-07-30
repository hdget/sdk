// Package qcc 企查查开放平台 API SDK
// 文档: https://openapi.qcc.com
package qcc

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/hdget/sdk/common/protobuf"
)

// 默认接口地址
const defaultBaseURL = "https://api.qichacha.com"

// Api 企查查开放平台 API 接口
type Api interface {
	// FuzzySearch 企业模糊搜索
	// 通过搜索关键字（如企业名、人名、产品名、地址、电话、经营范围等）获取匹配的企业，
	// 每次请求最多返回 5 条记录
	// searchKey: 搜索关键词
	// pageIndex: 页码，可选，默认第 1 页
	FuzzySearch(searchKey string, list ...*protobuf.ListParam) (int64, []Enterprise, error)
}

type apiImpl struct {
	key       string // 应用 AppKey
	secretKey string // 密钥 SecretKey
	baseURL   string // 接口地址
}

// Option 配置选项
type Option func(*apiImpl)

// New 创建企查查 API 实例
// key 和 secretKey 均可在企查查开放平台官网（openapi.qcc.com）获取
func New(key, secretKey string, opts ...Option) Api {
	impl := &apiImpl{
		key:       key,
		secretKey: secretKey,
		baseURL:   defaultBaseURL,
	}
	for _, opt := range opts {
		opt(impl)
	}
	return impl
}

// genToken 生成鉴权 Token
// 规则: key + timespan + secretKey 拼接后做 32 位 md5 加密，取大写字符串
func genToken(key, timespan, secretKey string) string {
	sum := md5.Sum([]byte(key + timespan + secretKey))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// authHeaders 生成鉴权请求头
// Token: 加密签名值
// Timespan: 精确到秒的 Unix 时间戳
func (impl *apiImpl) authHeaders() map[string]string {
	timespan := strconv.FormatInt(time.Now().Unix(), 10)
	return map[string]string{
		"Token":    genToken(impl.key, timespan, impl.secretKey),
		"Timespan": timespan,
	}
}
