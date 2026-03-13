package kd100

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hdget/sdk/libs/logistics"
)

// API endpoints
const (
	InstantQueryURL = "https://poll.kuaidi100.com/poll/query.do"
	SubscribeURL    = "https://poll.kuaidi100.com/poll"
	RecognizeURL    = "https://www.kuaidi100.com/autonumber/auto"
)

// 供应商名称
const ProviderName = "kd100"

// init 注册快递100供应商
func init() {
	logistics.RegisterFactory(ProviderName, New)
}

// api 快递100 API实现
type api struct {
	appId      string
	appSecret  string
	httpClient *http.Client
}

// New 创建快递100 API实例
func New(cfg *logistics.Config) (logistics.LogisticsApi, error) {
	if cfg == nil || cfg.AppId == "" || cfg.AppSecret == "" {
		return nil, logistics.ErrEmptyConfig
	}

	return &api{
		appId:     cfg.AppId,
		appSecret: cfg.AppSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// sign 生成快递100签名
// sign = MD5(param + key + customer).toUpperCase()
func sign(param, key, customer string) string {
	h := md5.New()
	h.Write([]byte(param + key + customer))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// Query 即时查询物流轨迹
func (a *api) Query(ctx context.Context, req *logistics.QueryRequest) (*logistics.QueryResult, error) {
	if req.ShipperCode == "" {
		return nil, logistics.ErrInvalidShipperCode
	}
	if req.TrackingNo == "" {
		return nil, logistics.ErrInvalidTrackingNo
	}

	// 构建查询参数
	param := kd100QueryParam{
		Com:      req.ShipperCode,
		Num:      req.TrackingNo,
		Phone:    req.Phone,
		Resultv2: "4", // 开启行政区域解析
	}

	paramJSON, _ := json.Marshal(param)
	paramStr := string(paramJSON)

	// 构建请求
	formData := url.Values{}
	formData.Set("customer", a.appId)
	formData.Set("sign", sign(paramStr, a.appSecret, a.appId))
	formData.Set("param", paramStr)

	resp, err := a.httpClient.PostForm(InstantQueryURL, formData)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var kd100Resp kd100QueryResponse
	if err := json.Unmarshal(body, &kd100Resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if kd100Resp.Message != "" && kd100Resp.Message != "ok" {
		return nil, fmt.Errorf("%w: %s", logistics.ErrQueryFailed, kd100Resp.Message)
	}

	return &logistics.QueryResult{
		State:       convertStatus(kd100Resp.State),
		ShipperCode: kd100Resp.Com,
		TrackingNo:  kd100Resp.Nu,
		Traces:      convertTraces(kd100Resp.Data),
		Location:    kd100Resp.Location,
	}, nil
}

// Subscribe 订阅物流轨迹
func (a *api) Subscribe(ctx context.Context, req *logistics.SubscribeRequest) (*logistics.SubscribeResult, error) {
	if req.ShipperCode == "" {
		return nil, logistics.ErrInvalidShipperCode
	}
	if req.TrackingNo == "" {
		return nil, logistics.ErrInvalidTrackingNo
	}
	if req.CallbackURL == "" {
		return nil, fmt.Errorf("callback url is required")
	}

	// 构建订阅参数（包含租户ID）
	param := kd100SubscribeParam{
		Company: req.ShipperCode,
		Number:  req.TrackingNo,
		Key:     a.appSecret,
		Parameters: kd100SubscribeParameters{
			Callbackurl: req.CallbackURL,
			TID:         req.Tid, // 租户ID会原样返回
		},
	}

	paramJSON, _ := json.Marshal(param)
	paramStr := string(paramJSON)

	// 构建请求
	formData := url.Values{}
	formData.Set("schema", "json")
	formData.Set("param", paramStr)

	resp, err := a.httpClient.PostForm(SubscribeURL, formData)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var kd100Resp kd100SubscribeResponse
	if err := json.Unmarshal(body, &kd100Resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &logistics.SubscribeResult{
		Success: kd100Resp.ReturnCode == "200",
		Message: kd100Resp.Message,
	}, nil
}

// Recognize 识别快递公司
func (a *api) Recognize(ctx context.Context, trackingNo string) ([]logistics.RecognizeResult, error) {
	if trackingNo == "" {
		return nil, logistics.ErrInvalidTrackingNo
	}

	// 快递100单号识别是GET请求
	recognizeURL := fmt.Sprintf("%s?num=%s&key=%s", RecognizeURL, trackingNo, a.appSecret)

	resp, err := a.httpClient.Get(recognizeURL)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var kd100Resp []kd100RecognizeItem
	if err := json.Unmarshal(body, &kd100Resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	results := make([]logistics.RecognizeResult, 0, len(kd100Resp))
	for _, item := range kd100Resp {
		results = append(results, logistics.RecognizeResult{
			ShipperCode: item.ComCode,
			ShipperName: item.Name,
		})
	}

	return results, nil
}

// ParseCallback 解析回调数据
func (a *api) ParseCallback(data []byte) (*logistics.CallbackData, error) {
	var cb kd100Callback
	if err := json.Unmarshal(data, &cb); err != nil {
		return nil, fmt.Errorf("%w: %v", logistics.ErrParseCallbackFailed, err)
	}

	// 从Parameters中解析租户ID
	tenantID := cb.Parameters.TID

	return &logistics.CallbackData{
		ShipperCode: cb.Company,
		TrackingNo:  cb.Number,
		Tid:         tenantID,
		State:       convertStatus(cb.State),
		Traces:      convertTraces(cb.Data),
		Success:     cb.State == "3", // 签收状态
		CourierInfo: &logistics.CourierInfo{
			Name:  cb.CourierName,
			Phone: cb.CourierPhone,
		},
	}, nil
}

// BuildCallbackResponse 构建回调响应
func (a *api) BuildCallbackResponse(success bool, message string) []byte {
	resp := kd100CallbackResponse{
		Result:  success,
		Message: message,
	}
	data, _ := json.Marshal(resp)
	return data
}
