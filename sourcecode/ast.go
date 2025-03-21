package sourcecode

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type CallSignature struct {
	importPath    string
	functionChain string
}

// astIsFunctionCall 检查某个函数调用是导入包名和函数调用链条是否完全匹配
func astIsFunctionCall(n *ast.CallExpr, importMap map[string]string, fnCall *CallSignature) bool {
	if caller, ok := astGetCaller(n); ok {
		if importMap[caller] == fnCall.importPath {
			if astGetFunctionChain(n) == fnCall.functionChain {
				return true
			}
		}
	}
	return false
}

// astStructTypeBelongTo 检查第一个field是否是匿名引入的模块， e,g: type A struct { dapr.InvocationModule }
// if len(structures.Fields.List) > 0 {
// possibleModuleExpr := fmt.Sprintf("%s", structures.Fields.List[0].Type)
// if moduleName := moduleExpr2moduleName[possibleModuleExpr]; moduleName != "" {
// found, _ := filepath.Rel(srcDir, filepath.Dir(path))
// m.ModulePaths[moduleName] = filepath.ToSlash(found)
// }
// }
func astStructTypeBelongTo(n *ast.GenDecl, expr2name map[string]string) string {
	// 仅处理类型声明
	if n.Tok == token.TYPE {
		for _, spec := range n.Specs {
			// 如果类型规范是类型别名或类型声明
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				// 如果类型规范是结构体类型
				if structures, ok := typeSpec.Type.(*ast.StructType); ok {
					for _, field := range structures.Fields.List {
						fieldTypeExpr := fmt.Sprintf("%s", field.Type)
						if v, exists := expr2name[fieldTypeExpr]; exists {
							return v
						}
					}
				}
			}
		}
	}
	return ""
}

// astGetFunctionChain 获取完整的函数调用链
func astGetFunctionChain(n *ast.CallExpr) string {
	functions := astParseFunction(n)
	return strings.Join(pie.Reverse(functions), ".")
}

// parseMetaData 递归解析链式函数调用，最近的Ident.Name作为包名，最先调用的函数在slice的最前面
func astParseFunction(n *ast.CallExpr) []string {
	var methods []string

	// 递归提取方法名
	for {
		// 检查 Fun 是否是 SelectorExpr
		selectorExpr, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}

		// 添加方法名
		methods = append(methods, selectorExpr.Sel.Name)

		// 检查 X 是否是另一个 CallExpr
		nextCallExpr, ok := selectorExpr.X.(*ast.CallExpr)
		if !ok {
			break
		}

		// 继续递归
		n = nextCallExpr
	}
	return methods
}

// astGetImportPaths 获取导入路径
func astGetImportPaths(f *ast.File) map[string]string {
	importMap := make(map[string]string)
	for _, imp := range f.Imports {
		pkgName := ""
		if imp.Name != nil {
			pkgName = imp.Name.Name // 处理别名导入，如 `import alias "math/rand"`
		} else {
			// 提取完整路径（去掉引号）
			pkgPath := strings.Trim(imp.Path.Value, `"`)
			// 获取包名（路径的最后一部分）
			pkgName = pkgPath[strings.LastIndex(pkgPath, "/")+1:]
		}
		importMap[pkgName] = strings.Trim(imp.Path.Value, `"`)
	}
	return importMap
}

// 获取调用者
func astGetCaller(n *ast.CallExpr) (string, bool) {
	// 递归查找调用者
	for {
		// 检查 Fun 是否是 SelectorExpr
		selectorExpr, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}

		// 检查 X 是否是另一个 CallExpr
		nextCallExpr, ok := selectorExpr.X.(*ast.CallExpr)
		if !ok {
			// 如果不是 CallExpr，则可能是调用者（如 sdk）
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				return ident.Name, true
			}
			break
		}

		// 继续递归
		n = nextCallExpr
	}

	return "", false
}

//
//// getFunctionChain 获取完整的函数调用链
//func gastGetFunctionChain(n *ast.CallExpr) string {
//	functions := make([]string, 0)
//	astRecursiveParseFunction(n, functions)
//	return strings.Join(pie.Reverse(functions), ".")
//}

// astGetEmbedInfo 获取嵌入资源的信息，返回变量名，embed路径
func astGetEmbedVarAndRelPath(n *ast.GenDecl) (string, string, bool) {
	// 如果是 GenDecl 类型，则可能是 import 或者变量声明等
	if n.Tok == token.VAR {
		for _, spec := range n.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				if astIsEmbedFSType(valueSpec.Type) {
					return valueSpec.Names[0].Name, astGetEmbedRelPath(n), true
				}
			}
		}
	}
	return "", "", false
}

// 检查类型是否为 embed.FS
func astIsEmbedFSType(expr ast.Expr) bool {
	if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
		if ident, ok := selectorExpr.X.(*ast.Ident); ok && ident.Name == "embed" {
			if selectorExpr.Sel.Name == "FS" {
				return true
			}
		}
	}
	return false
}

// 获取 embed 路径
func astGetEmbedRelPath(n *ast.GenDecl) string {
	// 如果直接定义变量
	// //go:embed assets/*
	// var assets embed.FS
	if n.Doc != nil {
		for _, comment := range n.Doc.List {
			if strings.HasPrefix(comment.Text, "//go:embed") {
				// 提取路径部分
				return filepath.Dir(strings.TrimSpace(strings.TrimPrefix(comment.Text, "//go:embed")))
			}
		}
	}
	// 如果定义在var block中
	// var (
	//   //go:embed assets/*
	//   assets embed.FS
	// )
	for _, spec := range n.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			if valueSpec.Doc != nil {
				for _, comment := range valueSpec.Doc.List {
					if strings.HasPrefix(comment.Text, "//go:embed") {
						// 提取路径部分
						return filepath.Dir(strings.TrimSpace(strings.TrimPrefix(comment.Text, "//go:embed")))
					}
				}
			}
		}
	}
	return ""
}

// Parse 尝试从源代码中查找嵌入路径, 返回嵌入资源的绝对路径和相对路径
func astParseEmbed(callerFilePath string) (string, string, error) {
	// 创建一个新的文件集
	fset := token.NewFileSet()

	// 解析源文件，同时保留注释
	f, err := parser.ParseFile(fset, callerFilePath, nil, parser.ParseComments)
	if err != nil {
		return "", "", err
	}

	// 遍历AST节点
	count := 0
	var foundVar, foundRelPath, embedAbsPath string
	ast.Inspect(f, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.GenDecl:
			if varName, relPath, ok := astGetEmbedVarAndRelPath(n); ok {
				foundVar = varName
				foundRelPath = relPath
				return false
			}
		}
		count += 1
		return foundVar == ""
	})

	fmt.Println("xxxxxxxxxxxxxx:", count)

	if foundVar == "" {
		return "", "", fmt.Errorf("embed.FS variable declare not found, var: %s", foundVar)
	}

	// 有可能定义了embed.FS,但是没有指定编译指令//go:embed
	if foundRelPath == "" {
		return "", "", fmt.Errorf("//go:embed compiler directive not found, var: %s", foundVar)
	}

	if foundRelPath == "." {
		return "", "", fmt.Errorf("//go:embed must specify a directory, var: %s", foundVar)
	}

	embedAbsPath = filepath.Join(filepath.Dir(callerFilePath), foundRelPath)
	return embedAbsPath, foundRelPath, nil
}
