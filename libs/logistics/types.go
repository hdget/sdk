package logistics

// QueryRequest 即时查询请求
type QueryRequest struct {
	ShipperCode string `json:"shipper_code"` // 快递公司编码
	TrackingNo  string `json:"tracking_no"`  // 快递单号
	ExtraInfo   string `json:"extra_info"`   // 额外信息（顺丰/中通:手机号, 京东:商家编码, 其他按需传递）
}

// QueryResult 即时查询结果
type QueryResult struct {
	State         State     `json:"state"`                    // 物流状态
	ShipperCode   string    `json:"shipper_code"`             // 快递公司编码
	TrackingNo    string    `json:"tracking_no"`              // 快递单号
	Traces        []Trace   `json:"traces"`                   // 物流轨迹
	Location      string    `json:"location"`                 // 当前位置
	EstimatedTime string    `json:"estimated_time,omitempty"` // 预计到达时间
	CourierInfo   *CourierInfo `json:"courier_info,omitempty"` // 快递员信息
}

// Trace 物流轨迹
type Trace struct {
	Time     string `json:"time"`     // 时间
	Content  string `json:"content"`  // 内容
	Location string `json:"location"` // 当前位置
	Status   State  `json:"status"`   // 轨迹状态
}

// CourierInfo 快递员信息
type CourierInfo struct {
	Name  string `json:"name"`  // 快递员姓名
	Phone string `json:"phone"` // 快递员电话
}

// PickUpInfo 取件信息（入柜/驿站时返回）
type PickUpInfo struct {
	Code    string `json:"code"`    // 取件码
	Address string `json:"address"` // 取件地址
	Station string `json:"station"` // 取件站点名称
}

// SubscribeRequest 订阅请求
type SubscribeRequest struct {
	ShipperCode string   `json:"shipper_code"`              // 快递公司编码
	TrackingNo  string   `json:"tracking_no"`               // 快递单号
	ExtraInfo   string   `json:"extra_info"`                // 额外信息（顺丰/中通:手机号, 京东:商家编码, 其他按需传递）
	CallbackURL string   `json:"callback_url"`              // 回调地址
	Metadata    string   `json:"metadata"`                  // 回调元数据，回调时会原封不动带回来
	Sender      *Contact `json:"sender,omitempty"`          // 发件人信息
	Receiver    *Contact `json:"receiver,omitempty"`        // 收件人信息
}

// SubscribeResult 订阅结果
type SubscribeResult struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
}

// Contact 联系人信息
type Contact struct {
	Name     string `json:"name"`     // 姓名
	Phone    string `json:"phone"`    // 电话
	Province string `json:"province"` // 省
	City     string `json:"city"`     // 市
	District string `json:"district"` // 区
	Address  string `json:"address"`  // 详细地址
}

// RecognizeResult 识别结果
type RecognizeResult struct {
	ShipperCode string `json:"shipper_code"` // 快递公司编码
	ShipperName string `json:"shipper_name"` // 快递公司名称
}

// CallbackData 回调数据
type CallbackData struct {
	ShipperCode string       `json:"shipper_code"`             // 快递公司编码
	TrackingNo  string       `json:"tracking_no"`              // 快递单号
	MetaData    string       `json:"metadata"`                 // 元数据,订阅时带过去，回调会带回来
	Status      Status       `json:"status"`                   // 状态信息（统一状态 + 原始状态码 + 描述）
	Traces      []Trace      `json:"traces"`                   // 轨迹列表
	Location    string       `json:"location"`                 // 当前位置/所在城市
	Success     bool         `json:"success"`                  // 是否成功
	Reason      string       `json:"reason"`                   // 失败原因
	CourierInfo *CourierInfo `json:"courier_info,omitempty"`   // 快递员信息
	PickUpInfo  *PickUpInfo  `json:"pick_up_info,omitempty"`   // 取件信息（入柜/驿站时返回）
}

// CallbackResponse 回调响应
type CallbackResponse struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
}

// ApiVendor API供应商
type ApiVendor string

const (
	VendorKdniao ApiVendor = "kdniao" // 快递鸟
	VendorKd100  ApiVendor = "kd100"  // 快递100
)

// Config 物流API统一配置
type Config struct {
	Name      string `mapstructure:"name"`       // 供应商名称: kdniao, kd100
	AppId     string `mapstructure:"app_id"`     // 应用ID
	AppSecret string `mapstructure:"app_secret"` // 应用密钥
}
