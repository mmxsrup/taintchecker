package taintchecker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "taintchecker is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "taintchecker",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	inspect.Preorder(nil, func(n ast.Node) {
		var id *ast.Ident
		var args []ast.Expr;

		switch n := n.(type) {
		case *ast.CallExpr:
			args = n.Args
			switch fun := n.Fun.(type) {
			case *ast.Ident:
				id = fun
			case *ast.SelectorExpr:
				id = fun.Sel
			}
		}

		flag := false

		if id != nil && id.Name == "ReadFile" {
			for _, arg := range args {
				if ident, ok := arg.(*ast.Ident); ok {
					flag = checkTaintNode(ident)
				}
			}
		}

		if flag {
			pass.Reportf(n.Pos(), "NG")
		}
	})

	return nil, nil
}

func checkTaintNode(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		if node.Obj == nil || node.Obj.Kind == ast.Con {
			return true
		}
		if node, ok := node.Obj.Decl.(ast.Node); ok {
			return checkTaintNode(node)
		}
		return false
	case *ast.AssignStmt:
		if len(node.Rhs) == 0 {
			return false
		}
		for _, rval := range node.Rhs {
			if checkTaintNode(rval) {
				return false
			}
		}
		return true
	}
	return false
}
