package types

import "context"

// LogisticsStatus 物流状态抽象
type LogisticsStatus int

const (
	LogisticsStatusPickedUp  LogisticsStatus = 1 // 已揽收
	LogisticsStatusInTransit LogisticsStatus = 2 // 运输中
	LogisticsStatusException LogisticsStatus = 3 // 异常
	LogisticsStatusSigned    LogisticsStatus = 4 // 已签收
)

// Trace 物流轨迹
type Trace struct {
	Action        string // 动作代码
	AcceptTime    string // 时间
	AcceptStation string // 描述
	Location      string // 地点
	Remark        string // 备注
}

// TraceItem 轨迹数据项
type TraceItem struct {
	ShipperCode    string          // 快递公司编码
	LogisticCode   string          // 快递单号
	Status         LogisticsStatus // 抽象状态
	Success        bool            // 是否成功
	Reason         string          // 失败原因
	Callback       string          // 回传字段（含target标识）
	Traces         []Trace         // 轨迹列表
	DeliveryMan    string          // 快递员
	DeliveryManTel string          // 快递员电话
}

// InstantQueryResult 即时查询结果
type InstantQueryResult struct {
	ShipperCode  string
	LogisticCode string
	Status       LogisticsStatus
	Traces       []Trace
}

// RecognizeResult 单号识别结果
type RecognizeResult struct {
	LogisticCode string
	Shippers     []ShipperInfo
}

// ShipperInfo 快递公司信息
type ShipperInfo struct {
	Code string // 快递公司编码
	Name string // 快递公司名称
}

// Contact 联系人信息
type Contact struct {
	Name         string
	Mobile       string
	ProvinceName string
	CityName     string
	ExpAreaName  string
	Address      string
}

// LogisticsProvider 物流服务Provider接口
type LogisticsProvider interface {
	My() LogisticsClient            // 获取默认客户端
	By(name string) LogisticsClient // 获取指定名称的客户端
}

// LogisticsClient 物流服务客户端接口
type LogisticsClient interface {
	// InstantQuery 即时查询物流轨迹
	InstantQuery(ctx context.Context, shipperCode, logisticCode string) (*InstantQueryResult, error)

	// InstantQueryWithCustomer 即时查询（带客户名，SF必填）
	InstantQueryWithCustomer(ctx context.Context, shipperCode, logisticCode, customerName string) (*InstantQueryResult, error)

	// Subscribe 订阅物流轨迹
	Subscribe(ctx context.Context, shipperCode, logisticCode, callback string) error

	// SubscribeWithOptions 订阅物流轨迹（带选项）
	SubscribeWithOptions(ctx context.Context, shipperCode, logisticCode, callback string, sender, receiver *Contact) error

	// Recognize 识别快递单号所属快递公司
	Recognize(ctx context.Context, logisticCode string) (*RecognizeResult, error)
}

// LogisticsPushHandler 物流推送处理器接口
type LogisticsPushHandler interface {
	Handle(ctx context.Context, item TraceItem) error
}

// NamedLogisticsPushHandler 带名称的处理器
type NamedLogisticsPushHandler struct {
	Target  string               // 目标标识
	Handler LogisticsPushHandler // 处理器
}

// LogisticsPushHandlerFunc 函数适配器
type LogisticsPushHandlerFunc func(ctx context.Context, item TraceItem) error

// Handle implements LogisticsPushHandler
func (f LogisticsPushHandlerFunc) Handle(ctx context.Context, item TraceItem) error {
	return f(ctx, item)
}
