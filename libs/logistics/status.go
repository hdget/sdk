package logistics

// Status 状态信息（统一状态 + 原始状态码 + 原始状态描述）
type Status struct {
	State State  `json:"state"` // 统一状态（用于业务逻辑判断）
	Code  string `json:"code"`  // 原始状态码（如 "201"、"301"）
	Desc  string `json:"desc"`  // 原始状态描述（如 "到达派件城市"、"本人签收"）
}

// State 物流状态（统一内部状态）
type State int

const (
	StateUnknown    State = iota // 未知
	StateNoTrace                 // 暂无轨迹
	StateCollected               // 已揽收
	StateInTransit               // 在途中
	StateDelivering              // 派件中
	StateSigned                  // 已签收
	StateProblem                 // 问题件
	StateReturned                // 退回
	StateRejected                // 拒签
	StateCleared                 // 清关中
)

// String 返回状态的字符串表示
func (s State) String() string {
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
func (s State) IsTerminal() bool {
	return s == StateSigned || s == StateRejected || s == StateReturned
}

// IsSuccess 判断是否为成功状态
func (s State) IsSuccess() bool {
	return s == StateSigned
}

// IsProblem 判断是否为问题状态
func (s State) IsProblem() bool {
	return s == StateProblem || s == StateRejected || s == StateReturned
}
