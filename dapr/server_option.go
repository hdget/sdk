package dapr

//type ServerOption func(impl *daprServerImpl)
//
//type Action func() error
//
//// WithServerEntry 提供服务器启动入口的信息
//func WithServerEntry(serverImportPath, serverRunFunction string) ServerOption {
//	return func(impl *daprServerImpl) {
//		impl.serverImportPath = serverImportPath
//		impl.serverRunFuncName = serverRunFunction
//	}
//}
//
//// WithBeforeRunActions 提供在服务器启动之前的操作
//func WithBeforeRunActions(actions ...Action) ServerOption {
//	return func(impl *daprServerImpl) {
//		impl.actions = append(impl.actions, actions...)
//	}
//}
