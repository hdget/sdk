package dapr

import (
	"github.com/elliotchance/pie/v2"
	"go/ast"
	"strings"
)

type astCallExprParser struct {
	pkg     string
	fnNames []string
}

func newAstCallExprParser(n *ast.CallExpr) *astCallExprParser {
	p := &astCallExprParser{
		pkg:     "",
		fnNames: make([]string, 0),
	}
	p.parse(n)
	return p
}

func (p *astCallExprParser) getFunctionName() string {
	return strings.Join(pie.Reverse(p.fnNames), ".")
}

// parse 递归解析链式函数调用，最近的Ident.Name作为包名，最先调用的函数在slice的最前面
func (p *astCallExprParser) parse(n *ast.CallExpr) {
	selExpr, ok := n.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	ident, ok := selExpr.X.(*ast.Ident)
	if ok {
		p.pkg = ident.Name
		p.fnNames = append(p.fnNames, selExpr.Sel.Name)
		return
	}

	// 有可能链式函数调用，需要递归检查
	innerCallExpr, isInnerCall := selExpr.X.(*ast.CallExpr)
	if isInnerCall {
		p.fnNames = append(p.fnNames, selExpr.Sel.Name)
		p.parse(innerCallExpr)
	}
}
