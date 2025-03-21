package sourcecode

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// 通过ast解析导入的包，以及调用server.New().Run来判断服务器运行的入口包
type parseMetaHandler struct {
	sc *sourceCodeManagerImpl
}

//// astEmbedPathResult 嵌入资源的路径
//type astEmbedPathResult struct {
//	varName string // 变量名字
//	absPath string
//	relPath string
//}

var (
	moduleExpr2moduleName = map[string]string{
		"&{dapr InvocationModule}": "InvocationModule", // 服务调用模块
		"&{dapr EventModule}":      "EventModule",      // 事件模块
		"&{dapr HealthModule}":     "HealthModule",     // 健康检测模块
		"&{dapr DelayEventModule}": "DelayEventModule", // 延迟事件模块
	}
)

const (
	// 需要将invocation module的包导入到server包来保证dapr方法的自动注册
	defaultAppServerImportPath  = "github.com/hdget/sdk/dapr" // 服务缺省的导入的包路径
	defaultAppServerRunFunction = "NewGrpcServer"             // 函数签名：dapr.NewGrpcServer
)

func newParseMetaHandler(sc *sourceCodeManagerImpl) Handler {
	return &parseMetaHandler{
		sc: sc,
	}
}

// Handle 从代码中解析meta信息并保存起来
// 1. dapr module paths
// 2. server.Start的入口文件
// 3. embed.FS的绝对路径和相对路径
func (h *parseMetaHandler) Handle() error {
	fmt.Println("")
	fmt.Printf("=== parse source code meta data ===")
	fmt.Println("")

	meta, err := h.parseMetaData(h.sc.skipDirs...)
	if err != nil {
		return err
	}

	if err = newMetaDataManager(h.sc.srcDir).Store(meta); err != nil {
		return err
	}

	if err = newMetaDataManager(h.sc.srcDir).Print(); err != nil {
		return err
	}

	return nil
}

func (h *parseMetaHandler) parseMetaData(skipDirs ...string) (*metaData, error) {
	st, err := os.Stat(h.sc.srcDir)
	if err != nil {
		return nil, errors.Wrapf(err, "sourc code dir not accessable, srcDir: %s", h.sc.srcDir)
	}

	if !st.IsDir() {
		return nil, fmt.Errorf("invalid source code dir, dir: %s", h.sc.srcDir)
	}

	result := &metaData{
		ModulePaths: make(map[string]string),
	}

	// 遍历源代码
	_ = filepath.Walk(h.sc.srcDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// 忽略掉.开头的隐藏目录和需要skip的目录
			if strings.HasPrefix(info.Name(), ".") || pie.Contains(skipDirs, info.Name()) {
				return filepath.SkipDir
			}
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return errors.Wrapf(err, "golang ast parseMetaData file, path: %s", path)
			}

			// 新建一个记录导入别名与包名映射关系的字典
			pkgName2importPath := astGetImportPaths(f)

			ast.Inspect(f, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.GenDecl: // import, constant, type or variable declaration
					// 判断结构声明中的field类型是否出现在moduleExpr2moduleName中，如果出现了，说明类似下面的type声明找到了
					//
					//	type v1_example struct {
					//		dapr.InvocationModule
					//	}
					//
					if moduleName := astStructTypeBelongTo(n, moduleExpr2moduleName); moduleName != "" {
						relPath, _ := filepath.Rel(h.sc.srcDir, filepath.Dir(path))
						result.ModulePaths[moduleName] = filepath.ToSlash(relPath)
					}
				case *ast.CallExpr: // 函数调用
					// 查找server.Start调用所在的文件路径
					if result.EntryPath == "" {
						if astIsFunctionCall(n, pkgName2importPath, &functionCall{
							importPath:    defaultAppServerImportPath,
							functionChain: defaultAppServerRunFunction,
						}) {
							result.EntryPath = path
						}
					}
				}
				return true
			})
		}

		return nil
	})

	// 通过遍历结果来获取serverEntry和相关数据
	return result, nil
}
