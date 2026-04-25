package kdniao

import "github.com/hdget/sdk/libs/logistics"

// instantQueryRequest 即时查询请求
type instantQueryRequest struct {
	ShipperCode  string `json:"ShipperCode"`            // 快递公司编码
	LogisticCode string `json:"LogisticCode"`           // 快递单号
	CustomerName string `json:"CustomerName,omitempty"` // 客户名称(SF必填)
}

// instantQueryResponse 即时查询响应
type instantQueryResponse struct {
	Success        bool    `json:"Success"`
	EBusinessID    string  `json:"EBusinessID"`
	ShipperCode    string  `json:"ShipperCode"`
	LogisticCode   string  `json:"LogisticCode"`
	State          string  `json:"State"`          // 物流状态
	StateEx        string  `json:"StateEx"`        // 增值物流状态
	Location       string  `json:"Location"`       // 所在城市
	Reason         string  `json:"Reason"`         // 失败原因
	Traces         []trace `json:"Traces"`         // 物流轨迹
	DeliveryMan    string  `json:"DeliveryMan"`    // 快递员
	DeliveryManTel string  `json:"DeliveryManTel"` // 快递员电话
}

// trace 快递鸟轨迹
type trace struct {
	Action        string `json:"Action"`        // 动作代码
	AcceptTime    string `json:"AcceptTime"`    // 时间
	AcceptStation string `json:"AcceptStation"` // 描述
	Location      string `json:"Location"`      // 地点
	Remark        string `json:"Remark"`        // 备注
}

// subscribeRequest 轨迹订阅请求
type subscribeRequest struct {
	ShipperCode  string         `json:"ShipperCode"`            // 快递公司编码
	LogisticCode string         `json:"LogisticCode"`           // 快递单号
	CustomerName string         `json:"CustomerName,omitempty"` // 客户名称(SF:手机号后四位, JD:商家编码)
	Callback     string         `json:"Callback,omitempty"`     // 用户自定义回调字段（限50字符）
	Sender       *kdniaoContact `json:"Sender,omitempty"`       // 发件人
	Receiver     *kdniaoContact `json:"Receiver,omitempty"`     // 收件人
}

// subscribeResponse 订阅响应
type subscribeResponse struct {
	Success     bool   `json:"Success"`
	EBusinessID string `json:"EBusinessID"`
	ResultCode  string `json:"ResultCode"`
	Reason      string `json:"Reason"`
}

// recognizeRequest 单号识别请求
type recognizeRequest struct {
	LogisticCode string `json:"LogisticCode"` // 快递单号
}

// recognizeResponse 单号识别响应
type recognizeResponse struct {
	Success      bool            `json:"Success"`
	EBusinessID  string          `json:"EBusinessID"`
	LogisticCode string          `json:"LogisticCode"`
	Shippers     []kdniaoShipper `json:"Shippers"`
}

// kdniaoShipper 快递公司信息
type kdniaoShipper struct {
	ShipperCode string `json:"ShipperCode"` // 快递公司编码
	ShipperName string `json:"ShipperName"` // 快递公司名称
}

// pushRequest 快递鸟推送请求（RequestData中的实际数据）
type pushRequest struct {
	EBusinessID string         `json:"EBusinessID"`
	Count       string         `json:"Count"`
	PushTime    string         `json:"PushTime"`
	Data        []pushDataItem `json:"Data"`
}

// pushDataItem 推送数据项（根据文档4.2.2.5完整字段）
type pushDataItem struct {
	EBusinessID    string       `json:"EBusinessID"`
	ShipperCode    string       `json:"ShipperCode"`
	LogisticCode   string       `json:"LogisticCode"`
	Callback       string       `json:"Callback"`       // 用户自定义回调字段
	Success        bool         `json:"Success"`        // 是否成功
	Reason         string       `json:"Reason"`         // 失败原因
	State          string       `json:"State"`          // 物流状态
	StateEx        string       `json:"StateEx"`        // 增值物流状态
	Location       string       `json:"Location"`       // 所在城市
	Traces         []trace      `json:"Traces"`         // 物流轨迹
	DeliveryMan    string       `json:"DeliveryMan"`    // 快递员
	DeliveryManTel string       `json:"DeliveryManTel"` // 快递员电话
	PickUpInfo     *pickUpInfo  `json:"PickUpInfo"`     // 取件信息
	NextCity       string       `json:"NextCity"`       // 下一站城市
}

// pickUpInfo 取件信息（入柜/驿站时返回）
type pickUpInfo struct {
	PickUpCode    string `json:"PickUpCode"`    // 取件码
	PickUpAddress string `json:"PickUpAddress"` // 取件地址
	PickUpStation string `json:"PickUpStation"` // 取件站点名称
}

// pushResponse 推送响应（根据文档4.2.2.6）
type pushResponse struct {
	EBusinessID string `json:"EBusinessID"` // 用户ID
	UpdateTime  string `json:"UpdateTime"`  // 更新时间
	Success     bool   `json:"Success"`     // 是否成功
	Reason      string `json:"Reason"`      // 失败原因
}

// kdniaoContact 快递鸟联系人
type kdniaoContact struct {
	Name         string `json:"Name,omitempty"`
	Mobile       string `json:"Mobile,omitempty"`
	ProvinceName string `json:"ProvinceName,omitempty"`
	CityName     string `json:"CityName,omitempty"`
	AreaName     string `json:"AreaName,omitempty"`
	Address      string `json:"Address,omitempty"`
}

// convertContact 转换联系人
func convertContact(c *logistics.Contact) *kdniaoContact {
	if c == nil {
		return nil
	}
	return &kdniaoContact{
		Name:         c.Name,
		Mobile:       c.Phone,
		ProvinceName: c.Province,
		CityName:     c.City,
		AreaName:     c.District,
		Address:      c.Address,
	}
}

// callbackForm 快递鸟回调表单数据（外层参数）
type callbackForm struct {
	RequestType string `schema:"RequestType"` // 101=普通订阅推送, 102=增值订阅推送
	EBusinessID string `schema:"EBusinessID"`
	RequestData string `schema:"RequestData"` // JSON 格式的轨迹数据
	DataSign    string `schema:"DataSign"`    // 签名
	DataType    string `schema:"DataType"`    // 数据类型，2=JSON
}
