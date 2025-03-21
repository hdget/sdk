package sourcecode

type Handler interface {
	Handle() error
}

type HandlerFunc func(*sourceCodeManagerImpl) Handler

var (
	patchHandlers = []HandlerFunc{
		newParseMetaHandler,            // 解析源代码数据
		newAddModuleImportPathsHandler, // 添加模块导入路径
	}

	inspectHandlers = []HandlerFunc{
		newFindRouteAnnotationHandler, // 查找路由注解
	}
)
