package kdniao

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hdget/sdk/libs/logistics"
)

// API endpoints
const (
	InstantQueryURL = "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
	SubscribeURL    = "https://api.kdniao.com/api/dist"
	RecognizeURL    = "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
)

// Request types
const (
	RequestTypeInstantQuery = "1002"
	RequestTypeSubscribe    = "8008"
	RequestTypeRecognize    = "2002"
)

// DataType for KDNiao API (2 = JSON)
const DataType = "2"

// 供应商名称
const ProviderName = "kdniao"

// init 注册快递鸟供应商
func init() {
	logistics.RegisterFactory(ProviderName, New)
}

// api 快递鸟API实现
type api struct {
	appId      string
	appSecret  string
	httpClient *http.Client
}

// New 创建快递鸟API实例
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

// doRequest 执行HTTP请求
func (a *api) doRequest(apiURL string, requestType string, requestData interface{}) ([]byte, error) {
	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("marshal request data: %w", err)
	}
	requestDataStr := string(jsonData)

	// 生成签名
	dataSign := sign(requestDataStr, a.appSecret)

	// 构建请求参数
	formData := url.Values{}
	formData.Set("RequestData", requestDataStr)
	formData.Set("EBusinessID", a.appId)
	formData.Set("RequestType", requestType)
	formData.Set("DataSign", dataSign)
	formData.Set("DataType", DataType)

	// 发送POST请求
	resp, err := a.httpClient.PostForm(apiURL, formData)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	return body, nil
}

// Query 即时查询物流轨迹
func (a *api) Query(ctx context.Context, req *logistics.QueryRequest) (*logistics.QueryResult, error) {
	if req.ShipperCode == "" {
		return nil, logistics.ErrInvalidShipperCode
	}
	if req.TrackingNo == "" {
		return nil, logistics.ErrInvalidTrackingNo
	}

	kdniaoReq := &instantQueryRequest{
		ShipperCode:  req.ShipperCode,
		LogisticCode: req.TrackingNo,
		CustomerName: req.Phone,
	}

	var resp instantQueryResponse
	err := a.doRequestWithResponse(InstantQueryURL, RequestTypeInstantQuery, kdniaoReq, &resp)
	if err != nil {
		return nil, fmt.Errorf("instant query: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("%w: %s", logistics.ErrQueryFailed, resp.Reason)
	}

	return &logistics.QueryResult{
		State:       convertStatus(resp.State),
		ShipperCode: resp.ShipperCode,
		TrackingNo:  resp.LogisticCode,
		Traces:      convertTraces(resp.Traces),
		Location:    resp.Location,
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
	// 注意：快递鸟的回调地址在官网后台配置，不在订阅请求中传递

	kdniaoReq := &subscribeRequest{
		ShipperCode:  req.ShipperCode,
		LogisticCode: req.TrackingNo,
		Callback:     req.Tid, // 使用 Callback 字段传递租户ID（限50字符）
		Sender:       convertContact(req.Sender),
		Receiver:     convertContact(req.Receiver),
	}

	var resp subscribeResponse
	err := a.doRequestWithResponse(SubscribeURL, RequestTypeSubscribe, kdniaoReq, &resp)
	if err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}

	return &logistics.SubscribeResult{
		Success: resp.Success,
		Message: resp.Reason,
	}, nil
}

// Recognize 识别快递公司
func (a *api) Recognize(ctx context.Context, trackingNo string) ([]logistics.RecognizeResult, error) {
	if trackingNo == "" {
		return nil, logistics.ErrInvalidTrackingNo
	}

	kdniaoReq := &recognizeRequest{
		LogisticCode: trackingNo,
	}

	var resp recognizeResponse
	err := a.doRequestWithResponse(RecognizeURL, RequestTypeRecognize, kdniaoReq, &resp)
	if err != nil {
		return nil, fmt.Errorf("recognize: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("%w: unrecognized tracking number", logistics.ErrRecognizeFailed)
	}

	results := make([]logistics.RecognizeResult, 0, len(resp.Shippers))
	for _, s := range resp.Shippers {
		results = append(results, logistics.RecognizeResult{
			ShipperCode: s.ShipperCode,
			ShipperName: s.ShipperName,
		})
	}

	return results, nil
}

// ParseCallback 解析回调数据
func (a *api) ParseCallback(data []byte) (*logistics.CallbackData, error) {
	var req pushRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("%w: %v", logistics.ErrParseCallbackFailed, err)
	}

	if len(req.Data) == 0 {
		return nil, fmt.Errorf("%w: no data in callback", logistics.ErrParseCallbackFailed)
	}

	// 取第一条数据
	item := req.Data[0]

	return &logistics.CallbackData{
		ShipperCode: item.ShipperCode,
		TrackingNo:  item.LogisticCode,
		Tid:         item.Callback, // 从 Callback 获取租户ID
		State:       convertStatus(item.State),
		Traces:      convertTraces(item.Traces),
		Success:     item.Success,
		Reason:      item.Reason,
		CourierInfo: &logistics.CourierInfo{
			Name:  item.DeliveryMan,
			Phone: item.DeliveryManTel,
		},
	}, nil
}

// BuildCallbackResponse 构建回调响应
func (a *api) BuildCallbackResponse(success bool, message string) []byte {
	resp := pushResponse{
		Success: success,
		Reason:  message,
	}
	data, _ := json.Marshal(resp)
	return data
}

// doRequestWithResponse 执行HTTP请求并解析响应
func (a *api) doRequestWithResponse(apiURL string, requestType string, requestData interface{}, response interface{}) error {
	body, err := a.doRequest(apiURL, requestType, requestData)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("decode response: %w, body: %s", err, string(body))
	}

	return nil
}
