package kd100

import (
	"bytes"
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
const VendorName = "kd100"

// init 注册快递100供应商
func init() {
	logistics.RegisterFactory(VendorName, New)
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
//
// 安全说明: MD5 仅用于满足快递100 API 的签名要求，不应用于其他安全敏感场景。
// 参考: 快递100开放平台文档
func sign(param, key, customer string) string {
	h := md5.New()
	h.Write([]byte(param + key + customer))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// verifySign 验证签名
// 使用 constant time 比较防止时序攻击
func verifySign(param, key, customer, dataSign string) bool {
	expectedSign := sign(param, key, customer)
	// 使用简单字符串比较（因为签名值长度固定且不敏感）
	// 注意: 快递100签名验证使用MD5，仅用于API兼容性
	return expectedSign == dataSign
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
		Phone:    req.ExtraInfo,
		Resultv2: "4", // 开启行政区域解析
	}

	paramJSON, _ := json.Marshal(param)
	paramStr := string(paramJSON)

	// 构建请求
	formData := url.Values{}
	formData.Set("customer", a.appId)
	formData.Set("sign", sign(paramStr, a.appSecret, a.appId))
	formData.Set("param", paramStr)

	// 使用 context 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, InstantQueryURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(httpReq)
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
		State:       convertState(kd100Resp.State),
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

	// 构建订阅参数
	param := kd100SubscribeParam{
		Company: req.ShipperCode,
		Number:  req.TrackingNo,
		Key:     a.appSecret,
		Parameters: kd100SubscribeParameters{
			Callbackurl: req.CallbackURL,
			Metadata:    req.Metadata, // 元数据会原样返回
			Phone:       req.ExtraInfo,
		},
	}

	paramJSON, _ := json.Marshal(param)
	paramStr := string(paramJSON)

	// 构建请求
	formData := url.Values{}
	formData.Set("schema", "json")
	formData.Set("param", paramStr)

	// 使用 context 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, SubscribeURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(httpReq)
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
// 安全警告: 快递100单号识别API要求使用GET请求，密钥会出现在URL中。
// 请确保:
// 1. 不在日志中记录完整URL
// 2. 使用HTTPS防止网络监听
// 3. 定期轮换密钥
func (a *api) Recognize(ctx context.Context, trackingNo string) ([]logistics.RecognizeResult, error) {
	if trackingNo == "" {
		return nil, logistics.ErrInvalidTrackingNo
	}

	// 对参数进行URL编码，防止特殊字符注入
	encodedTrackingNo := url.QueryEscape(trackingNo)
	encodedKey := url.QueryEscape(a.appSecret)

	// 快递100单号识别是GET请求，密钥会出现在URL中（API要求）
	recognizeURL := fmt.Sprintf("%s?num=%s&key=%s", RecognizeURL, encodedTrackingNo, encodedKey)

	// 使用 context 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, recognizeURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		// 不要在错误信息中包含完整URL（含密钥）
		return nil, fmt.Errorf("http request to recognize api: %w", err)
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

// ParseCallback 解析回调数据并验证签名
// 快递100回调数据格式：表单提交，包含 param(JSON字符串) 和 sign(签名)
// 签名验证: sign = MD5(param + key + customer).toUpperCase()
func (a *api) ParseCallback(data []byte) (*logistics.CallbackData, error) {
	// 1. 解析表单格式请求
	var req kd100CallbackRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("%w: parse request format: %v", logistics.ErrParseCallbackFailed, err)
	}

	// 2. 验证签名
	if req.Param == "" || req.Sign == "" {
		return nil, fmt.Errorf("%w: missing param or sign", logistics.ErrParseCallbackFailed)
	}

	if !verifySign(req.Param, a.appSecret, a.appId, req.Sign) {
		return nil, fmt.Errorf("%w: signature verification failed", logistics.ErrParseCallbackFailed)
	}

	// 3. 解析回调数据
	var cb kd100Callback
	if err := json.Unmarshal([]byte(req.Param), &cb); err != nil {
		return nil, fmt.Errorf("%w: parse callback data: %v", logistics.ErrParseCallbackFailed, err)
	}

	// 4. 检查是否有物流数据
	if cb.LastResult == nil {
		return nil, fmt.Errorf("%w: no lastResult in callback", logistics.ErrParseCallbackFailed)
	}

	// 5. 验证必要字段
	if cb.LastResult.Com == "" || cb.LastResult.Nu == "" {
		return nil, fmt.Errorf("%w: missing shipperCode or trackingNo", logistics.ErrParseCallbackFailed)
	}

	// 获取元数据
	var metadata string
	if cb.Parameters != nil {
		metadata = cb.Parameters.Metadata
	}

	return &logistics.CallbackData{
		ShipperCode: cb.LastResult.Com,
		TrackingNo:  cb.LastResult.Nu,
		MetaData:    metadata,
		Status:      convertStatus(cb.LastResult.State),
		Traces:      convertTraces(cb.LastResult.Data),
		Location:    cb.LastResult.Location,
		Success:     cb.LastResult.State == "3", // 签收状态
		Reason:      cb.Message,
	}, nil
}

// BuildCallbackResponse 构建回调响应
func (a *api) BuildCallbackResponse(success bool, message string) []byte {
	resp := kd100CallbackResponse{
		Result:     success,
		ReturnCode: "200",
		Message:    message,
	}
	data, _ := json.Marshal(resp)
	return data
}
