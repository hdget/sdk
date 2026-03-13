package logistics

// LogisticsState 物流状态（统一内部状态）
type LogisticsState int

const (
	StateUnknown    LogisticsState = iota // 未知
	StateNoTrace                          // 暂无轨迹
	StateCollected                        // 已揽收
	StateInTransit                        // 在途中
	StateDelivering                       // 派件中
	StateSigned                           // 已签收
	StateProblem                          // 问题件
	StateReturned                         // 退回
	StateRejected                         // 拒签
	StateCleared                          // 清关中
)

// String 返回状态的字符串表示
func (s LogisticsState) String() string {
	switch s {
	case StateNoTrace:
		return "暂无轨迹"
	case StateCollected:
		return "已揽收"
	case StateInTransit:
		return "在途中"
	case StateDelivering:
		return "派件中"
	case StateSigned:
		return "已签收"
	case StateProblem:
		return "问题件"
	case StateReturned:
		return "退回"
	case StateRejected:
		return "拒签"
	case StateCleared:
		return "清关中"
	default:
		return "未知"
	}
}

// IsTerminal 判断是否为终态
func (s LogisticsState) IsTerminal() bool {
	return s == StateSigned || s == StateRejected || s == StateReturned
}

// IsSuccess 判断是否为成功状态
func (s LogisticsState) IsSuccess() bool {
	return s == StateSigned
}

// IsProblem 判断是否为问题状态
func (s LogisticsState) IsProblem() bool {
	return s == StateProblem || s == StateRejected || s == StateReturned
}