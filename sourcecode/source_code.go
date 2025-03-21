package sourcecode

import (
	"embed"
	"encoding/json"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/sdk/dapr"
	"runtime"
)

type SourceCodeManager interface {
	Patch() error                              // 处理源代码
	Inspect() error                            // 检查源代码
	GetRoutes() ([]*protobuf.RouteItem, error) // 获取路由
}

// SourceCodeInfo 源代码信息
type SourceCodeInfo struct {
	ModulePaths map[string]string // 模块的路径
	ServerEntry string            // appServer.Run的入口文件即appServer开始运行所在的go文件
}

type sourceCodeManagerImpl struct {
	assetManager        *assetManager
	metaManager         *metaDataManager
	srcDir              string
	skipDirs            []string
	serverImportPath    string // 服务器代码所在的包路径，即server.New().Run所在的包路径
	serverRunFuncName   string // 服务器代码运行的函数名
	handlerNameMatchers []dapr.HandlerNameMatcher
}

// New 初始化源代码管理器
func New(fs embed.FS, srcDir string, skipDirs ...string) (SourceCodeManager, error) {
	// 解析go:embed路径
	_, callerPath, _, _ := runtime.Caller(1)
	embedAbsPath, embedRelPath, err := astParseEmbed(callerPath)
	if err != nil {
		return nil, err
	}

	return &sourceCodeManagerImpl{
		assetManager:      newAssetManager(fs, embedAbsPath, embedRelPath),
		metaManager:       newMetaDataManager(srcDir),
		srcDir:            srcDir,
		skipDirs:          skipDirs,
		serverImportPath:  defaultAppServerImportPath,
		serverRunFuncName: defaultAppServerRunFunction,
	}, nil
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

func (impl *sourceCodeManagerImpl) GetRoutes() ([]*protobuf.RouteItem, error) {
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
