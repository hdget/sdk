package api

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hdget/sdk/libs/wechat/pkg/cache"
	"github.com/hdget/utils"
	"github.com/pkg/errors"
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
		return ret, errors.Wrapf(err, "http get request, url: %s, req: %v", url, req)
	}

	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return ret, errors.Wrapf(err, "parse result, url: %s, req: %v, ret: %s", url, req, utils.BytesToString(resp.Body()))
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
		return ret, errors.Wrapf(err, "http post request, url: %s, req: %v", url, req)
	}

	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return ret, errors.Wrapf(err, "parse result, url: %s, req: %v, ret: %s", url, req, utils.BytesToString(resp.Body()))
	}

	return ret, nil
}

// PostResponse http post and get response
func PostResponse(url string, request ...any) ([]byte, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	resp, err := resty.New().R().Post(url)
	if err != nil {
		return nil, errors.Wrapf(err, "http post request, url: %s, req: %v", url, req)
	}
	return resp.Body(), nil
}

func CheckResult(result Result, url string, request ...any) error {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("wx api error, url: %s, req: %v, ret: %v", url, req, result)
	}

	return nil
}

type RetrieveFunc func() (string, int, error) // 返回值，过期时间

func CacheFirst(objCache cache.ObjectCache, retrieveFunc RetrieveFunc) (string, error) {
	cached, _ := objCache.Get()
	if cached != "" {
		return cached, nil
	}

	value, expiresIn, err := retrieveFunc()
	if err != nil {
		return "", err
	}

	_ = objCache.Set(value, expiresIn)

	return value, nil
}
