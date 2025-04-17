package dapr

// NewInvocationModule 新建服务调用模块会执行下列操作:
// 1. 实例化invocation module
// 2. 注册invocation functions
// 3. 注册module
//func NewInvocationModule(app string, moduleObject InvocationModule, functions map[string]InvocationFunction) error {
//	// 首先实例化module
//	module, err := asInvocationModule(app, moduleObject)
//	if err != nil {
//		return err
//	}
//
//	// 然后注册handlers
//	err = module.RegisterHandlers(functions)
//	if err != nil {
//		return err
//	}
//
//	// 最后注册module
//	registerModule(module)
//
//	return nil
//}
//
//// DiscoverHandlers 获取Module作为receiver的所有MethodMatchFunction匹配的方法, MethodMatchFunction生成新的方法名和判断是否匹配
//func (m *invocationModuleImpl) discoverHandlers(args ...HandlerMatcher) ([]invocationHandler, error) {
//	matchFn := m.defaultHandlerNameMatcher
//	if len(args) > 0 {
//		matchFn = args[0]
//	}
//
//	handlers := make([]invocationHandler, 0)
//	// 这里需要传入当前实际正在使用的服务模块，即带有common.ServiceInvocationHandler的struct实例
//	for methodName, method := range reflectUtils.MatchReceiverMethods(m.self, InvocationFunction(nil)) {
//		handlerName, matched := matchFn(methodName)
//		if !matched {
//			continue
//		}
//
//		fn, err := m.toInvocationFunction(method)
//		if err != nil {
//			return nil, err
//		}
//
//		handlers = append(handlers, m.newInvocationHandler(m.Module, handlerName, fn))
//	}
//
//	return handlers, nil
//}

//func (m *invocationModuleImpl) toInvocationFunction(fn any) (InvocationFunction, error) {
//	realFunction, ok := fn.(InvocationFunction)
//	// 如果不是DaprInvocationHandler, 可能为实际的函数体
//	if !ok {
//		realFunction, ok = fn.(func(context.Context, *common.InvocationEvent) (any, error))
//		if !ok {
//			return nil, errInvalidInvocationFunction
//		}
//	}
//	return realFunction, nil
//}

//
//// matchHandlerSuffix 匹配方法名是否以handler结尾并将新方法名转为SnakeCase格式
//func (m *invocationModuleImpl) defaultHandlerNameMatcher(methodName string) (string, bool) {
//	lastIndex := strings.LastIndex(strings.ToLower(methodName), strings.ToLower(handlerNameSuffix))
//	if lastIndex <= 0 {
//		return "", false
//	}
//	return methodName[:lastIndex], true
//}
