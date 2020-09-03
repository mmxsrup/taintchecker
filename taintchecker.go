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
				switch a := arg.(type) {
				case *ast.Ident:
					if (*a.Obj).Kind == ast.Var {
						flag = true
					}
				}
			}
		}

		if flag {
			pass.Reportf(n.Pos(), "NG")
		}
	})

	return nil, nil
}
