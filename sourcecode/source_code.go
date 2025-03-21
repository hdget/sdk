package sourcecode

import (
	"embed"
	"encoding/json"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/sdk/dapr"
	"runtime"
)

type SourceCodeManager interface {
	Patch() error                                        // 处理源代码
	Inspect() error                                      // 检查源代码
	GetRouteAnnotations() ([]*protobuf.RouteItem, error) // 获取路由注解
}

type sourceCodeManagerImpl struct {
	assetManager        *assetManager
	metaManager         *metaDataManager
	srcDir              string
	skipDirs            []string
	serverCallSignature *CallSignature        // 服务调用签名, 用来判断server.Start在哪个文件执行，需要定位改文件来加入module import路径
	handlerMatchers     []dapr.HandlerMatcher // dapr module handler匹配规则
}

// New 初始化源代码管理器
func New(fs embed.FS, srcDir string, options ...Option) (SourceCodeManager, error) {
	// 解析go:embed路径
	_, callerPath, _, _ := runtime.Caller(1)
	embedAbsPath, embedRelPath, err := astParseEmbed(callerPath)
	if err != nil {
		return nil, err
	}

	m := &sourceCodeManagerImpl{
		assetManager: newAssetManager(fs, embedAbsPath, embedRelPath),
		metaManager:  newMetaDataManager(srcDir),
		srcDir:       srcDir,
		skipDirs:     make([]string, 0),
		serverCallSignature: &CallSignature{
			importPath:    defaultServerImportPath,
			functionChain: defaultServerNewFunction,
		},
	}

	for _, apply := range options {
		apply(m)
	}

	return m, nil
}

func (impl *sourceCodeManagerImpl) Patch() error {
	// call patch handlers
	for _, fn := range patchHandlers {
		if err := fn(impl).Handle(); err != nil {
			return err
		}
	}
	return nil
}

func (impl *sourceCodeManagerImpl) Inspect() error {
	// call inspect handlers
	for _, fn := range inspectHandlers {
		if err := fn(impl).Handle(); err != nil {
			return err
		}
	}

	return nil
}

func (impl *sourceCodeManagerImpl) GetRouteAnnotations() ([]*protobuf.RouteItem, error) {
	data, err := impl.assetManager.Load(fileRoutes)
	if err != nil {
		return nil, err
	}

	var routes []*protobuf.RouteItem
	err = json.Unmarshal(data, &routes)
	if err != nil {
		return nil, err
	}
	return routes, nil
}
