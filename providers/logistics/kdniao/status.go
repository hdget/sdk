package logistics_kdniao

import (
	"github.com/hdget/sdk/common/types"
)

// convertKdniaoStatus 将快递鸟状态转换为抽象状态
func convertKdniaoStatus(kdniaoStatus int) types.LogisticsStatus {
	switch kdniaoStatus {
	case StatusPickedUp:
		return types.LogisticsStatusPickedUp
	case StatusInTransit, StatusDelivering, StatusSending:
		return types.LogisticsStatusInTransit
	case StatusException, StatusProblem, StatusRejected:
		return types.LogisticsStatusException
	case StatusSigned, StatusPartialDelivery:
		return types.LogisticsStatusSigned
	default:
		return types.LogisticsStatusInTransit
	}
}

// convertTraces 将快递鸟轨迹转换为抽象轨迹
func convertTraces(kdniaoTraces []KdniaoTrace) []types.Trace {
	traces := make([]types.Trace, 0, len(kdniaoTraces))
	for _, t := range kdniaoTraces {
		traces = append(traces, types.Trace{
			Action:        t.Action,
			AcceptTime:    t.AcceptTime,
			AcceptStation: t.AcceptStation,
			Location:      t.Location,
			Remark:        t.Remark,
		})
	}
	return traces
}
