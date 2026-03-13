package kd100

import "github.com/hdget/sdk/libs/logistics"

// kd100QueryParam 即时查询参数
type kd100QueryParam struct {
	Com      string `json:"com"`                // 快递公司编码
	Num      string `json:"num"`                // 快递单号
	Phone    string `json:"phone,omitempty"`    // 电话（顺丰必填）
	Resultv2 string `json:"resultv2,omitempty"` // 扩展参数
}

// kd100QueryResponse 即时查询响应
type kd100QueryResponse struct {
	Message  string           `json:"message"`  // 错误消息
	State    string           `json:"state"`    // 状态码
	Status   string           `json:"status"`   // 当前状态
	Com      string           `json:"com"`      // 快递公司编码
	Nu       string           `json:"nu"`       // 快递单号
	Data     []kd100TraceData `json:"data"`     // 轨迹数据
	Location string           `json:"location"` // 当前位置
}

// kd100TraceData 轨迹数据
type kd100TraceData struct {
	Time     string `json:"time"`     // 时间
	Ftime    string `json:"ftime"`    // 格式化时间
	Context  string `json:"context"`  // 内容
	Location string `json:"location"` // 位置
	Status   string `json:"status"`   // 状态
}

// kd100SubscribeParam 订阅参数
type kd100SubscribeParam struct {
	Company    string                    `json:"company"`
	Number     string                    `json:"number"`
	Key        string                    `json:"key"`
	Parameters kd100SubscribeParameters  `json:"parameters"`
}

// kd100SubscribeParameters 订阅参数扩展
type kd100SubscribeParameters struct {
	Callbackurl string `json:"callbackurl"` // 回调地址
	TID         string `json:"tid,omitempty"` // 租户ID（回调时原样返回）
}

// kd100SubscribeResponse 订阅响应
type kd100SubscribeResponse struct {
	ReturnCode string `json:"returnCode"` // 返回码
	Message    string `json:"message"`    // 消息
}

// kd100RecognizeItem 识别结果项
type kd100RecognizeItem struct {
	ComCode string `json:"comCode"` // 快递公司编码
	Name    string `json:"name"`    // 快递公司名称 (需额外请求获取)
}

// kd100Callback 回调数据
type kd100Callback struct {
	Company      string                  `json:"company"`       // 快递公司
	Number       string                  `json:"number"`        // 快递单号
	State        string                  `json:"state"`         // 状态
	Status       string                  `json:"status"`        // 当前状态
	Data         []kd100TraceData        `json:"data"`          // 轨迹数据
	CourierName  string                  `json:"courierName"`   // 快递员姓名
	CourierPhone string                  `json:"courierPhone"`  // 快递员电话
	Parameters   kd100CallbackParameters `json:"parameters"`    // 订阅时传递的参数（原样返回）
}

// kd100CallbackParameters 回调参数（订阅时传递，回调时原样返回）
type kd100CallbackParameters struct {
	TID string `json:"tid"` // 租户ID
}

// kd100CallbackResponse 回调响应
type kd100CallbackResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

// convertStatus 将快递100状态转换为统一状态
func convertStatus(kd100State string) logistics.LogisticsState {
	switch kd100State {
	case "0": // 在途中
		return logistics.StateInTransit
	case "1": // 已揽收
		return logistics.StateCollected
	case "2": // 疑难
		return logistics.StateProblem
	case "3": // 已签收
		return logistics.StateSigned
	case "4": // 退签
		return logistics.StateReturned
	case "5": // 派件中
		return logistics.StateDelivering
	case "6": // 退回
		return logistics.StateReturned
	case "7": // 转投
		return logistics.StateInTransit
	case "8": // 清关
		return logistics.StateCleared
	case "14": // 拒签
		return logistics.StateRejected
	default:
		return logistics.StateUnknown
	}
}

// convertTraces 将快递100轨迹转换为统一轨迹
func convertTraces(kd100Traces []kd100TraceData) []logistics.Trace {
	traces := make([]logistics.Trace, 0, len(kd100Traces))
	for _, t := range kd100Traces {
		traces = append(traces, logistics.Trace{
			Time:     t.Ftime,
			Content:  t.Context,
			Location: t.Location,
		})
	}
	return traces
}