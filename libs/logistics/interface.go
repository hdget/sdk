package logistics

import "context"

// LogisticsApi 物流服务API接口
type LogisticsApi interface {
	// Query 即时查询物流轨迹
	Query(ctx context.Context, req *QueryRequest) (*QueryResult, error)

	// Subscribe 订阅物流轨迹
	Subscribe(ctx context.Context, req *SubscribeRequest) (*SubscribeResult, error)

	// Recognize 识别快递公司
	Recognize(ctx context.Context, trackingNo string) ([]RecognizeResult, error)

	// ParseCallback 解析回调数据
	ParseCallback(data []byte) (*CallbackData, error)

	// BuildCallbackResponse 构建回调响应
	BuildCallbackResponse(success bool, message string) []byte

	// GetShipperCode 根据快递公司名称获取编码
	GetShipperCode(name string) string
}