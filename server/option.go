package server

type Option func(impl *appServerImpl)

// WithServerEntry 提供服务器启动入口的信息
func WithServerEntry(serverImportPath, serverRunFunction string) Option {
	return func(impl *appServerImpl) {
		impl.serverImportPath = serverImportPath
		impl.serverRunFuncName = serverRunFunction
	}
}

// WithBeforeRunActions 提供在服务器启动之前的操作
func WithBeforeRunActions(actions ...Action) Option {
	return func(impl *appServerImpl) {
		impl.actions = append(impl.actions, actions...)
	}
}
