package devops

type Option func(impl *devOpsImpl)

func WithTableOperator(items ...TableOperator) Option {
	return func(impl *devOpsImpl) {
		impl.tableOperators = items
	}
}
