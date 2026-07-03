package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hdget/sdk/libs/wechat/pkg/cache"
	"github.com/hdget/utils"
	"github.com/pkg/errors"
)

// 敏感参数列表，不应出现在日志中
var sensitiveParams = []string{
	"secret", "appsecret", "access_token", "refresh_token",
	"password", "pwd", "key", "token", "session_key",
}

const (
	networkTimeout = 3 * time.Second
)

func Get[RESULT any](url string, request ...any) (RESULT, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	var ret RESULT
	resp, err := resty.New().SetTimeout(networkTimeout).
		SetHeader("Content-Type", "application/json; charset=UTF-8").
		R().SetBody(req).Get(url)
	if err != nil {
		return ret, errors.Wrapf(err, "http get request, url: %s", sanitizeURL(url))
	}

	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return ret, errors.Wrapf(err, "parse result, url: %s, ret: %s", sanitizeURL(url), utils.BytesToString(resp.Body()))
	}

	return ret, nil
}

func Post[RESULT any](url string, request ...any) (RESULT, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	var ret RESULT
	resp, err := resty.New().SetTimeout(networkTimeout).
		SetHeader("Content-Type", "application/json; charset=UTF-8").
		R().SetBody(req).Post(url)
	if err != nil {
		return ret, errors.Wrapf(err, "http post request, url: %s", sanitizeURL(url))
	}

	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return ret, errors.Wrapf(err, "parse result, url: %s, ret: %s", sanitizeURL(url), utils.BytesToString(resp.Body()))
	}

	return ret, nil
}

// PostResponse http post and get response
// 注意: 已添加超时设置
func PostResponse(url string, request ...any) ([]byte, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	resp, err := resty.New().SetTimeout(networkTimeout).R().SetBody(req).Post(url)
	if err != nil {
		return nil, errors.Wrapf(err, "http post request, url: %s", sanitizeURL(url))
	}
	return resp.Body(), nil
}

func CheckResult(result Result, url string, request ...any) error {
	if result.ErrCode != 0 {
		return fmt.Errorf("wx api error, url: %s, ret: %v", sanitizeURL(url), result)
	}

	return nil
}

type RetrieveFunc func() (string, int, error) // 返回值，过期时间

// CacheFirst 缓存优先策略：先从缓存获取，失败则从源获取并缓存
// 注意: 缓存操作错误会被记录但不会阻断流程，保证服务的可用性
func CacheFirst(objCache cache.ObjectCache, retrieveFunc RetrieveFunc) (string, error) {
	// 尝试从缓存获取
	cached, cacheErr := objCache.Get()
	if cacheErr == nil && cached != "" {
		return cached, nil
	}
	// 缓存获取失败或为空，继续从源获取（不阻断流程）

	// 从源获取数据
	value, expiresIn, err := retrieveFunc()
	if err != nil {
		return "", err
	}

	// 尝试设置缓存，失败不阻断流程
	setErr := objCache.Set(value, expiresIn)
	if setErr != nil {
		// 缓存设置失败不影响返回结果，但可以选择记录日志
		// 生产环境中应该添加日志记录：log.Warn("cache set failed", "err", setErr)
	}

	return value, nil
}

// sanitizeURL 过滤URL中的敏感参数值
// 将敏感参数的值替换为 "***"
func sanitizeURL(url string) string {
	result := url
	for _, param := range sensitiveParams {
		// 替换 URL 参数形式: param=value
		result = sanitizeQueryParam(result, param)
	}
	return result
}

// sanitizeQueryParam 替换特定参数的值
func sanitizeQueryParam(url string, paramName string) string {
	// 查找参数位置
	lowerURL := strings.ToLower(url)
	lowerParam := strings.ToLower(paramName)

	// 尝试找到参数: paramName= 或 &paramName=
	idx := strings.Index(lowerURL, lowerParam+"=")
	if idx == -1 {
		return url
	}

	// 找到值的开始位置
	valueStart := idx + len(paramName) + 1
	if valueStart >= len(url) {
		return url
	}

	// 找到值的结束位置（下一个&或字符串结束）
	valueEnd := strings.Index(url[valueStart:], "&")
	if valueEnd == -1 {
		valueEnd = len(url)
	} else {
		valueEnd = valueStart + valueEnd
	}

	// 替换值为 "***"
	return url[:valueStart] + "***" + url[valueEnd:]
}
