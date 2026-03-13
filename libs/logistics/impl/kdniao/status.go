package kdniao

import "github.com/hdget/sdk/libs/logistics"

// 快递鸟状态码
const (
	statusNoRecord        = 0 // 无轨迹
	statusPickedUp        = 1 // 已揽收
	statusInTransit       = 2 // 在途中
	statusSigned          = 3 // 已签收
	statusProblem         = 4 // 问题件
	statusDelivering      = 5 // 转寄
	statusSending         = 6 // 清关
	statusException       = 7 // 异常
	statusRejected        = 8 // 拒收
	statusPartialDelivery = 9 // 部分签收
)

// convertStatus 将快递鸟状态转换为统一状态
func convertStatus(kdniaoStatus int) logistics.LogisticsState {
	switch kdniaoStatus {
	case statusNoRecord:
		return logistics.StateNoTrace
	case statusPickedUp:
		return logistics.StateCollected
	case statusInTransit, statusDelivering, statusSending:
		return logistics.StateInTransit
	case statusSigned, statusPartialDelivery:
		return logistics.StateSigned
	case statusProblem, statusException:
		return logistics.StateProblem
	case statusRejected:
		return logistics.StateRejected
	default:
		return logistics.StateUnknown
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