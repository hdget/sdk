package logistics_kdniao

// API endpoints
const (
	// 即时查询API
	InstantQueryURL = "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
	// 轨迹订阅API
	SubscribeURL = "https://api.kdniao.com/api/dist"
	// 单号识别API
	RecognizeURL = "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
)

// Request types for KDNiao API
const (
	// RequestTypeInstantQuery 即时查询
	RequestTypeInstantQuery = "1002"
	// RequestTypeSubscribe 轨迹订阅
	RequestTypeSubscribe = "8008"
	// RequestTypeRecognize 单号识别
	RequestTypeRecognize = "2002"
)

// DataType for KDNiao API (2 = JSON)
const DataType = "2"

// KDNiao status codes
const (
	StatusNoRecord        = 0  // 无轨迹
	StatusPickedUp        = 1  // 已揽收
	StatusInTransit       = 2  // 在途中
	StatusSigned          = 3  // 已签收
	StatusProblem         = 4  // 问题件
	StatusDelivering      = 5  // 转寄
	StatusSending         = 6  // 清关
	StatusException       = 7  // 异常
	StatusRejected        = 8  // 拒收
	StatusPartialDelivery = 9  // 部分签收
)

// CommonRequest 快递鸟公共请求参数
type CommonRequest struct {
	RequestData string `json:"RequestData"`
	EBusinessID string `json:"EBusinessID"`
	RequestType string `json:"RequestType"`
	DataSign    string `json:"DataSign"`
	DataType    string `json:"DataType"`
}

// InstantQueryRequest 即时查询请求
type InstantQueryRequest struct {
	ShipperCode  string `json:"ShipperCode"`           // 快递公司编码
	LogisticCode string `json:"LogisticCode"`          // 快递单号
	CustomerName string `json:"CustomerName,omitempty"` // 客户名称(SF必填)
}

// InstantQueryResponse 即时查询响应
type InstantQueryResponse struct {
	Success       bool              `json:"Success"`
	EBusinessID   string            `json:"EBusinessID"`
	ShipperCode   string            `json:"ShipperCode"`
	LogisticCode  string            `json:"LogisticCode"`
	State         int               `json:"State"`         // 物流状态:0-无轨迹,1-已揽收,2-在途中,3-已签收,4-问题件,5-转寄,6-清关,7-异常
	StateEx       int               `json:"StateEx"`       // 增值物流状态
	Location      string            `json:"Location"`      // 所在城市
	Reason        string            `json:"Reason"`        // 失败原因
	Traces        []KdniaoTrace     `json:"Traces"`        // 物流轨迹
	DeliveryMan   string            `json:"DeliveryMan"`   // 快递员
	DeliveryManTel string           `json:"DeliveryManTel"` // 快递员电话
}

// KdniaoTrace 快递鸟轨迹
type KdniaoTrace struct {
	Action        string `json:"Action"`        // 动作代码
	AcceptTime    string `json:"AcceptTime"`    // 时间
	AcceptStation string `json:"AcceptStation"` // 描述
	Location      string `json:"Location"`      // 地点
	Remark        string `json:"Remark"`        // 备注
}

// SubscribeRequest 轨迹订阅请求
type SubscribeRequest struct {
	ShipperCode  string         `json:"ShipperCode"`           // 快递公司编码
	LogisticCode string         `json:"LogisticCode"`          // 快递单号
	Callback     string         `json:"Callback"`              // 回调URL(含target标识)
	Sender       *KdniaoContact `json:"Sender,omitempty"`     // 发件人
	Receiver     *KdniaoContact `json:"Receiver,omitempty"`   // 收件人
}

// KdniaoContact 快递鸟联系人
type KdniaoContact struct {
	Name         string `json:"Name"`
	Mobile       string `json:"Mobile"`
	ProvinceName string `json:"ProvinceName"`
	CityName     string `json:"CityName"`
	ExpAreaName  string `json:"ExpAreaName"`
	Address      string `json:"Address"`
}

// SubscribeResponse 订阅响应
type SubscribeResponse struct {
	Success    bool   `json:"Success"`
	EBusinessID string `json:"EBusinessID"`
	ResultCode  string `json:"ResultCode"`
	Reason     string `json:"Reason"`
}

// RecognizeRequest 单号识别请求
type RecognizeRequest struct {
	LogisticCode string `json:"LogisticCode"` // 快递单号
}

// RecognizeResponse 单号识别响应
type RecognizeResponse struct {
	Success     bool               `json:"Success"`
	EBusinessID string             `json:"EBusinessID"`
	LogisticCode string            `json:"LogisticCode"`
	Shippers    []KdniaoShipper   `json:"Shippers"`
}

// KdniaoShipper 快递公司信息
type KdniaoShipper struct {
	ShipperCode string `json:"ShipperCode"` // 快递公司编码
	ShipperName string `json:"ShipperName"` // 快递公司名称
}

// PushRequest 快递鸟推送请求
type PushRequest struct {
	EBusinessID string           `json:"EBusinessID"`
	Count       int              `json:"Count"`
	PushTime    string           `json:"PushTime"`
	Data        []PushDataItem   `json:"Data"`
}

// PushDataItem 推送数据项
type PushDataItem struct {
	EBusinessID    string        `json:"EBusinessID"`
	ShipperCode    string        `json:"ShipperCode"`
	LogisticCode   string        `json:"LogisticCode"`
	Success        bool          `json:"Success"`
	Reason         string        `json:"Reason"`
	State          int           `json:"State"`
	Callback       string        `json:"Callback"`
	Traces         []KdniaoTrace `json:"Traces"`
	DeliveryMan    string        `json:"DeliveryMan"`
	DeliveryManTel string        `json:"DeliveryManTel"`
}

// PushResponse 推送响应
type PushResponse struct {
	Success bool   `json:"Success"`
	Reason  string `json:"Reason"`
}