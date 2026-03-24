package kdniao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	"github.com/hdget/sdk/libs/logistics"
	"github.com/hdget/utils"
)

// API endpoints
const (
	InstantQueryURL = "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
	SubscribeURL    = "https://api.kdniao.com/api/dist"
	RecognizeURL    = "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
)

// Request types
const (
	RequestTypeInstantQuery    = "1002" // 即时查询
	RequestTypeSubscribe       = "8008" // 订阅请求
	RequestTypeRecognize       = "2002" // 单号识别
	RequestTypeCallbackNormal  = "101"  // 普通订阅推送回调
	RequestTypeCallbackPremium = "102"  // 增值订阅推送回调
)

// DataType for KDNiao API (2 = JSON)
const DataType = "2"

// 供应商名称
const VendorName = "kdniao"

// schemaDecoder 表单解码器（线程安全）
var schemaDecoder = schema.NewDecoder()

// init 注册快递鸟供应商
func init() {
	logistics.RegisterFactory(VendorName, New)
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
func (a *api) doRequest(ctx context.Context, apiURL string, requestType string, requestData interface{}) ([]byte, error) {
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

	// 使用 context 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp, err := a.httpClient.Do(req)
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
		CustomerName: req.ExtraInfo,
	}

	var resp instantQueryResponse
	err := a.doRequestWithResponse(ctx, InstantQueryURL, RequestTypeInstantQuery, kdniaoReq, &resp)
	if err != nil {
		return nil, fmt.Errorf("instant query: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("%w: %s", logistics.ErrQueryFailed, resp.Reason)
	}

	return &logistics.QueryResult{
		State:       convertState(resp.State),
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
		CustomerName: req.ExtraInfo,
		Callback:     req.Metadata, // 使用 Callback 字段传递元数据
		Sender:       convertContact(req.Sender),
		Receiver:     convertContact(req.Receiver),
	}

	var resp subscribeResponse
	err := a.doRequestWithResponse(ctx, SubscribeURL, RequestTypeSubscribe, kdniaoReq, &resp)
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
	err := a.doRequestWithResponse(ctx, RecognizeURL, RequestTypeRecognize, kdniaoReq, &resp)
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
// 快递鸟回调数据格式: RequestType=101&EBusinessID=xxx&RequestData={JSON数据}&DataSign=xxx&DataType=2
func (a *api) ParseCallback(data []byte) (*logistics.CallbackData, error) {
	// 1. 解析 form-encoded 数据到结构体
	values, err := url.ParseQuery(utils.BytesToString(data))
	if err != nil {
		return nil, fmt.Errorf("%w: parse form data: %v", logistics.ErrParseCallbackFailed, err)
	}

	// 设置忽略未知字段, 如果不加这一行，Decode 会返回错误："schema: invalid path
	schemaDecoder.IgnoreUnknownKeys(true)
	var form callbackForm
	if err = schemaDecoder.Decode(&form, values); err != nil {
		return nil, fmt.Errorf("%w: decode callback form: %v", logistics.ErrParseCallbackFailed, err)
	}

	if form.RequestData == "" {
		return nil, fmt.Errorf("%w: empty RequestData", logistics.ErrParseCallbackFailed)
	}

	// RequestData 可能被 URL 编码，需要先解码
	requestData, err := url.QueryUnescape(form.RequestData)
	if err != nil {
		requestData = form.RequestData // 解码失败时使用原始值
	}

	// 2. 验证签名
	// 注意：回调的DataSign可能经过URL编码，需要先解码再验证
	// 根据文档2.2.1：推送接口RequestType为101/102不需要进行URL编码
	// 但实际推送时可能已编码，所以尝试两种方式验证
	decodedSign, err := url.QueryUnescape(form.DataSign)
	if err != nil {
		decodedSign = form.DataSign // 解码失败时使用原始值
	}

	// 不管怎么样尝试用解码后的签名验证
	if !verifySign(requestData, a.appSecret, decodedSign) {
		return nil, fmt.Errorf("%w: signature verification failed", logistics.ErrParseCallbackFailed)
	}

	// 3. 解析RequestData中的轨迹数据
	var req pushRequest
	if err = json.Unmarshal([]byte(requestData), &req); err != nil {
		return nil, fmt.Errorf("%w: parse RequestData: %v", logistics.ErrParseCallbackFailed, err)
	}

	if len(req.Data) == 0 {
		return nil, fmt.Errorf("%w: no data in callback", logistics.ErrParseCallbackFailed)
	}

	// 取第一条数据
	item := req.Data[0]

	// 构建取件信息
	var pkInfo *logistics.PickUpInfo
	if item.PickUpInfo != nil {
		pkInfo = &logistics.PickUpInfo{
			Code:    item.PickUpInfo.PickUpCode,
			Address: item.PickUpInfo.PickUpAddress,
			Station: item.PickUpInfo.PickUpStation,
		}
	}

	return &logistics.CallbackData{
		ShipperCode: item.ShipperCode,
		TrackingNo:  item.LogisticCode,
		MetaData:    item.Callback, // 快递鸟Callback就是自定义数据，订阅时候传入的时候可以回调带回来
		Status:      convertStatus(item.State, item.StateEx, item.Success, item.Reason),
		Traces:      convertTraces(item.Traces),
		Success:     item.Success,
		Reason:      item.Reason,
		Location:    item.Location,
		CourierInfo: &logistics.CourierInfo{
			Name:  item.DeliveryMan,
			Phone: item.DeliveryManTel,
		},
		PickUpInfo: pkInfo,
	}, nil
}

// BuildCallbackResponse 构建回调响应（根据文档4.2.2.6）
func (a *api) BuildCallbackResponse(success bool, message string) []byte {
	resp := pushResponse{
		EBusinessID: a.appId,
		UpdateTime:  time.Now().Format(time.DateTime),
		Success:     success,
		Reason:      message,
	}
	data, _ := json.Marshal(resp)
	return data
}

// doRequestWithResponse 执行HTTP请求并解析响应
func (a *api) doRequestWithResponse(ctx context.Context, apiURL string, requestType string, requestData any, response any) error {
	body, err := a.doRequest(ctx, apiURL, requestType, requestData)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, response); err != nil {
		// 不在错误中暴露响应体，避免敏感信息泄露
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
