package kdniao

import "github.com/hdget/sdk/libs/logistics"

// stateExInfo 快递鸟增值状态码映射表
var stateExInfo = map[string]struct {
	state logistics.State
	desc  string
}{
	"0": {logistics.StateUnknown, "暂无轨迹信息"},
	// 已揽收
	"1": {logistics.StateCollected, "已揽收"},
	// 在途/派件相关
	"2":   {logistics.StateInTransit, "在途中"},
	"201": {logistics.StateDelivering, "到达派件城市"},
	"202": {logistics.StateDelivering, "派件中"},
	"211": {logistics.StateDelivering, "已放入快递柜或驿站"},
	// 签收相关
	"3":   {logistics.StateSigned, "已签收"},
	"301": {logistics.StateSigned, "正常签收"},
	"302": {logistics.StateSigned, "派件异常后最终签收"},
	"304": {logistics.StateSigned, "代收签收"},
	"311": {logistics.StateSigned, "快递柜或驿站签收"},
	// 问题件相关
	"4":   {logistics.StateProblem, "问题件"},
	"401": {logistics.StateProblem, "发货无信息"},
	"402": {logistics.StateProblem, "超时未签收"},
	"403": {logistics.StateProblem, "超时未更新"},
	"404": {logistics.StateRejected, "拒收(退件)"},
	"405": {logistics.StateProblem, "派件异常"},
	"406": {logistics.StateSigned, "退货签收"}, // 退货签收，虽然首数字是4，但也属于正常轨迹
	"407": {logistics.StateProblem, "退货未签收"},
	"412": {logistics.StateSigned, "快递柜或驿站超时未取或已签收"}, // 入驿站超过18h的单，很可能其实已经被取走了，只是驿站不更新，一般情况下也可以视为签收
	"413": {logistics.StateProblem, "单号已拦截"},
	"10":  {logistics.StateCollected, "待揽件"},
}

// stateInfo 快递鸟基础状态映射表
var stateInfo = map[string]struct {
	state logistics.State
	desc  string
}{
	"0": {logistics.StateNoTrace, "无轨迹"},
	"1": {logistics.StateCollected, "已揽收"},
	"2": {logistics.StateInTransit, "在途中"},
	"3": {logistics.StateSigned, "已签收"},
	"4": {logistics.StateProblem, "问题件"},
}

// convertState 将快递鸟状态转换为统一状态
func convertState(state string) logistics.State {
	if v, ok := stateInfo[state]; ok {
		return v.state
	}
	return logistics.StateUnknown
}

// convertStatus 综合使用 State 和 StateEx 转换状态
// 优先使用 StateEx，StateEx 为空时使用 State
// 同时处理物流异常情况：Success=false 且 Reason 为"三天无轨迹"或"七天内无轨迹变化"
func convertStatus(state, stateEx string, success bool, reason string) logistics.Status {
	// 处理物流异常情况：三天无轨迹或七天内无轨迹变化
	if !success && (reason == "三天无轨迹" || reason == "七天内无轨迹变化") {
		return logistics.Status{
			State: logistics.StateProblem,
			Code:  "",
			Desc:  reason,
		}
	}

	// 优先使用 StateEx
	if stateEx != "" {
		if v, ok := stateExInfo[stateEx]; ok {
			return logistics.Status{
				State: v.state,
				Code:  stateEx,
				Desc:  v.desc,
			}
		}
	}
	// 回退到 State
	if v, ok := stateInfo[state]; ok {
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

// convertTraces 将快递鸟轨迹转换为统一轨迹
func convertTraces(kdniaoTraces []trace) []logistics.Trace {
	traces := make([]logistics.Trace, 0, len(kdniaoTraces))
	for _, t := range kdniaoTraces {
		traces = append(traces, logistics.Trace{
			Time:     t.AcceptTime,
			Content:  t.AcceptStation,
			Location: t.Location,
		})
	}
	return traces
}
