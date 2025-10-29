package devops

type Option func(impl *devOpsImpl)

func WithTableOperator(items ...TableOperator) Option {
	return func(impl *devOpsImpl) {
		impl.tableOperators = items
	}
}

// WithDangerConfirm 危险命令是否需要确认
func WithDangerConfirm(noConfirm bool) Option {
	return func(impl *devOpsImpl) {
		impl.needDangerConfirm = noConfirm
	}
}
