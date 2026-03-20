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
	Company    string                   `json:"company"`
	Number     string                   `json:"number"`
	Key        string                   `json:"key"`
	Parameters kd100SubscribeParameters `json:"parameters"`
}

// kd100SubscribeParameters 订阅参数扩展
type kd100SubscribeParameters struct {
	Callbackurl string `json:"callbackurl"`        // 回调地址
	Metadata    string `json:"metadata,omitempty"` // 元数据（回调时原样返回）
	Phone       string `json:"phone,omitempty"`    // 收寄件人电话（顺丰、中通必填）
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

// kd100Callback 回调数据（根据文档2.4）
type kd100Callback struct {
	Status     string                   `json:"status"`     // poll:监听状态
	BillStatus string                   `json:"billstatus"` // got:获取到快递信息
	Message    string                   `json:"message"`    // 错误信息
	AutoCheck  string                   `json:"autoCheck"`  // 是否自动判断公司
	ComOld     string                   `json:"comOld"`     // 原快递公司编码
	ComNew     string                   `json:"comNew"`     // 新快递公司编码
	LastResult *kd100LastResult         `json:"lastResult"` // 最新物流信息
	Parameters *kd100CallbackParameters `json:"parameters"` // 订阅时的回调参数
}

// kd100LastResult 最新物流信息（同即时查询返回）
type kd100LastResult struct {
	Message  string           `json:"message"`  // 消息
	State    string           `json:"state"`    // 状态码
	Status   string           `json:"status"`   // 当前状态
	Com      string           `json:"com"`      // 快递公司编码
	Nu       string           `json:"nu"`       // 快递单号
	Data     []kd100TraceData `json:"data"`     // 轨迹数据
	Location string           `json:"location"` // 当前位置
}

// kd100CallbackParameters 回调参数
type kd100CallbackParameters struct {
	Callbackurl string `json:"callbackurl"`        // 回调地址
	Metadata    string `json:"metadata,omitempty"` // 元数据（原样返回）
	Salt        string `json:"salt,omitempty"`     // 签名用随机字符串
}

// kd100CallbackResponse 回调响应
type kd100CallbackResponse struct {
	Result     bool   `json:"result"`
	ReturnCode string `json:"returnCode"`
	Message    string `json:"message"`
}

// kd100StateInfo 快递100状态映射表
var kd100StateInfo = map[string]struct {
	state logistics.State
	desc  string
}{
	"0":  {logistics.StateInTransit, "在途中"},
	"1":  {logistics.StateCollected, "已揽收"},
	"2":  {logistics.StateProblem, "疑难"},
	"3":  {logistics.StateSigned, "已签收"},
	"4":  {logistics.StateReturned, "退签"},
	"5":  {logistics.StateDelivering, "派件中"},
	"6":  {logistics.StateReturned, "退回"},
	"7":  {logistics.StateInTransit, "转投"},
	"8":  {logistics.StateCleared, "清关"},
	"14": {logistics.StateRejected, "拒签"},
}

// convertState 将快递100状态转换为统一状态
func convertState(state string) logistics.State {
	if v, ok := kd100StateInfo[state]; ok {
		return v.state
	}
	return logistics.StateUnknown
}

// convertStatus 将快递100状态转换为统一状态信息
// 返回：Status（统一状态 + 原始状态码 + 原始状态描述）
func convertStatus(state string) logistics.Status {
	if v, ok := kd100StateInfo[state]; ok {
		return logistics.Status{
			State: v.state,
			Code:  state,
			Desc:  v.desc,
		}
	}
	return logistics.Status{
		State: logistics.StateUnknown,
		Code:  state,
		Desc:  "未知",
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
