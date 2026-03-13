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
	Success        bool          `json:"Success"`
	EBusinessID    string        `json:"EBusinessID"`
	ShipperCode    string        `json:"ShipperCode"`
	LogisticCode   string        `json:"LogisticCode"`
	State          int           `json:"State"`          // 物流状态
	StateEx        int           `json:"StateEx"`        // 增值物流状态
	Location       string        `json:"Location"`       // 所在城市
	Reason         string        `json:"Reason"`         // 失败原因
	Traces         []trace `json:"Traces"`         // 物流轨迹
	DeliveryMan    string        `json:"DeliveryMan"`    // 快递员
	DeliveryManTel string        `json:"DeliveryManTel"` // 快递员电话
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
	ShipperCode  string         `json:"ShipperCode"`        // 快递公司编码
	LogisticCode string         `json:"LogisticCode"`       // 快递单号
	Callback     string         `json:"Callback,omitempty"` // 用户自定义回调字段（限50字符，用于传递租户ID）
	Sender       *kdniaoContact `json:"Sender,omitempty"`   // 发件人
	Receiver     *kdniaoContact `json:"Receiver,omitempty"` // 收件人
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

// pushRequest 快递鸟推送请求
type pushRequest struct {
	EBusinessID string         `json:"EBusinessID"`
	Count       int            `json:"Count"`
	PushTime    string         `json:"PushTime"`
	Data        []pushDataItem `json:"Data"`
}

// pushDataItem 推送数据项
type pushDataItem struct {
	EBusinessID    string  `json:"EBusinessID"`
	ShipperCode    string  `json:"ShipperCode"`
	LogisticCode   string  `json:"LogisticCode"`
	Callback       string  `json:"Callback"` // 用户自定义回调字段（订阅时传递的租户ID）
	Success        bool    `json:"Success"`
	Reason         string  `json:"Reason"`
	State          int     `json:"State"`
	Traces         []trace `json:"Traces"`
	DeliveryMan    string  `json:"DeliveryMan"`
	DeliveryManTel string  `json:"DeliveryManTel"`
}

// pushResponse 推送响应
type pushResponse struct {
	Success bool   `json:"Success"`
	Reason  string `json:"Reason"`
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