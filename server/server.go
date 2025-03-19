package server

import (
	"embed"
	"fmt"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/sdk"
	"github.com/hdget/sdk/dapr"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
)

type Action func() error

type AppServer interface {
	Run(appAddress string) error
	PatchSourceCode(srcDir string, skipDirs ...string) error // 处理源代码
	InspectSourceCode(srcDir string, skipDirs ...string) error
	GetRoutes() ([]*protobuf.RouteItem, error)
}

type appServerImpl struct {
	assetManager      AssetManager
	actions           []Action
	serverImportPath  string // 服务器代码所在的包路径，即server.New().Run所在的包路径
	serverRunFuncName string // 服务器代码运行的函数名
}

const (
	// 通过ast解析导入的包，以及调用server.New().Run来判断服务器运行的入口包
	// 需要将invocation module的包导入到server包来保证dapr方法的自动注册
	defaultAppServerImportPath  = "github.com/hdget/sdk/server" // 服务缺省的导入的包路径
	defaultAppServerRunFunction = "New.Run"                     // 函数签名：server.New().Run(appAddress)
)

func New(assetFs embed.FS, options ...Option) AppServer {
	_, callPath, _, _ := runtime.Caller(1)

	embedAbsPath, embedRelPath := findEmbedPath(callPath)

	srv := &appServerImpl{
		assetManager:      newAssetManager(embedAbsPath, embedRelPath, assetFs),
		actions:           make([]Action, 0),
		serverImportPath:  defaultAppServerImportPath,
		serverRunFuncName: defaultAppServerRunFunction,
	}

	for _, apply := range options {
		apply(srv)
	}
	return srv
}

func (impl *appServerImpl) Run(appAddress string) error {
	appServer, err := dapr.NewGrpcDaprServer(sdk.Logger(), appAddress)
	if err != nil {
		return errors.Wrap(err, "new app server")
	}

	if err = appServer.Start(); err != nil {
		sdk.Logger().Fatal("start app server", "err", err)
	}
	return nil
}

// 尝试查找查找嵌入路径, 返回绝对路径和相对路径
func findEmbedPath(filePath string) (string, string) {
	// 创建一个新的文件集
	fset := token.NewFileSet()

	// 解析源文件，同时保留注释
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", ""
	}

	// 遍历AST节点
	var absEmbedPath, relEmbedPath string
	ast.Inspect(file, func(n ast.Node) bool {
		// 如果是 GenDecl 类型，则可能是 import 或者变量声明等
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok && valueSpec.Doc != nil {
					// 检查每个注释行
					for _, comment := range valueSpec.Doc.List {
						if strings.HasPrefix(comment.Text, "//go:embed") {
							// 嵌入的相对路径
							relEmbedPath = filepath.Dir(strings.TrimSpace(strings.TrimPrefix(comment.Text, "//go:embed")))
							// 绝对路径
							absEmbedPath = filepath.Join(filepath.Dir(filePath), relEmbedPath)
							fmt.Printf("Found go:embed directive, file: %s, variable: %s, embedPath: %s\n", filePath, valueSpec.Names[0], relEmbedPath)
							break
						}
					}
				}
			}
		}

		// 找到embedPath后停止遍历
		if relEmbedPath != "" {
			return false
		}
		return true
	})

	return absEmbedPath, relEmbedPath
}
