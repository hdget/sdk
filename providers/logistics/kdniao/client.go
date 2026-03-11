package logistics_kdniao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hdget/sdk/common/types"
)

type kdniaoClient struct {
	eBusinessID string
	appKey      string
	httpClient  *http.Client
}

func newKdniaoClient(conf *kdniaoClientConfig) (types.LogisticsClient, error) {
	c := &kdniaoClient{
		eBusinessID: conf.EBusinessID,
		appKey:      conf.AppKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	return c, nil
}

// doRequest 执行HTTP请求
func (c *kdniaoClient) doRequest(apiURL string, requestType string, requestData interface{}) ([]byte, error) {
	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("marshal request data: %w", err)
	}
	requestDataStr := string(jsonData)

	// 生成签名
	dataSign := Sign(requestDataStr, c.appKey)

	// 构建请求参数
	formData := url.Values{}
	formData.Set("RequestData", requestDataStr)
	formData.Set("EBusinessID", c.eBusinessID)
	formData.Set("RequestType", requestType)
	formData.Set("DataSign", dataSign)
	formData.Set("DataType", DataType)

	// 发送POST请求
	resp, err := c.httpClient.PostForm(apiURL, formData)
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

// doRequestWithResponse 执行HTTP请求并解析响应
func (c *kdniaoClient) doRequestWithResponse(apiURL string, requestType string, requestData interface{}, response interface{}) error {
	body, err := c.doRequest(apiURL, requestType, requestData)
	if err != nil {
		return err
	}

	// 解析响应
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber() // 避免大数字精度丢失
	err = decoder.Decode(response)
	if err != nil {
		return fmt.Errorf("decode response: %w, body: %s", err, string(body))
	}

	return nil
}

// InstantQuery 即时查询物流轨迹
func (c *kdniaoClient) InstantQuery(ctx context.Context, shipperCode, logisticCode string) (*types.InstantQueryResult, error) {
	return c.InstantQueryWithCustomer(ctx, shipperCode, logisticCode, "")
}

// InstantQueryWithCustomer 即时查询（带客户名）
func (c *kdniaoClient) InstantQueryWithCustomer(ctx context.Context, shipperCode, logisticCode, customerName string) (*types.InstantQueryResult, error) {
	if shipperCode == "" {
		return nil, ErrInvalidShipperCode
	}
	if logisticCode == "" {
		return nil, ErrInvalidLogisticCode
	}

	req := &InstantQueryRequest{
		ShipperCode:  shipperCode,
		LogisticCode: logisticCode,
		CustomerName: customerName,
	}

	var resp InstantQueryResponse
	err := c.doRequestWithResponse(InstantQueryURL, RequestTypeInstantQuery, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("instant query: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("%w: %s", ErrQueryFailed, resp.Reason)
	}

	result := &types.InstantQueryResult{
		ShipperCode:  resp.ShipperCode,
		LogisticCode: resp.LogisticCode,
		Status:       convertKdniaoStatus(resp.State),
		Traces:       convertTraces(resp.Traces),
	}

	return result, nil
}

// Subscribe 订阅物流轨迹
func (c *kdniaoClient) Subscribe(ctx context.Context, shipperCode, logisticCode, callback string) error {
	return c.SubscribeWithOptions(ctx, shipperCode, logisticCode, callback, nil, nil)
}

// SubscribeWithOptions 订阅物流轨迹（带选项）
func (c *kdniaoClient) SubscribeWithOptions(ctx context.Context, shipperCode, logisticCode, callback string, sender, receiver *types.Contact) error {
	if shipperCode == "" {
		return ErrInvalidShipperCode
	}
	if logisticCode == "" {
		return ErrInvalidLogisticCode
	}
	if callback == "" {
		return fmt.Errorf("callback url is required")
	}

	req := &SubscribeRequest{
		ShipperCode:  shipperCode,
		LogisticCode: logisticCode,
		Callback:     callback,
	}

	if sender != nil {
		req.Sender = &KdniaoContact{
			Name:         sender.Name,
			Mobile:       sender.Mobile,
			ProvinceName: sender.ProvinceName,
			CityName:     sender.CityName,
			ExpAreaName:  sender.ExpAreaName,
			Address:      sender.Address,
		}
	}

	if receiver != nil {
		req.Receiver = &KdniaoContact{
			Name:         receiver.Name,
			Mobile:       receiver.Mobile,
			ProvinceName: receiver.ProvinceName,
			CityName:     receiver.CityName,
			ExpAreaName:  receiver.ExpAreaName,
			Address:      receiver.Address,
		}
	}

	var resp SubscribeResponse
	err := c.doRequestWithResponse(SubscribeURL, RequestTypeSubscribe, req, &resp)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("%w: %s", ErrSubscribeFailed, resp.Reason)
	}

	return nil
}

// Recognize 识别快递单号所属快递公司
func (c *kdniaoClient) Recognize(ctx context.Context, logisticCode string) (*types.RecognizeResult, error) {
	if logisticCode == "" {
		return nil, ErrInvalidLogisticCode
	}

	req := &RecognizeRequest{
		LogisticCode: logisticCode,
	}

	var resp RecognizeResponse
	err := c.doRequestWithResponse(RecognizeURL, RequestTypeRecognize, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("recognize: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("%w: unrecognized logistic code", ErrRecognizeFailed)
	}

	result := &types.RecognizeResult{
		LogisticCode: resp.LogisticCode,
		Shippers:     make([]types.ShipperInfo, 0, len(resp.Shippers)),
	}

	for _, s := range resp.Shippers {
		result.Shippers = append(result.Shippers, types.ShipperInfo{
			Code: s.ShipperCode,
			Name: s.ShipperName,
		})
	}

	return result, nil
}
