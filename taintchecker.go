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

var targetFunctions []string

func init() {
	// File access functions
	targetFunctions = append(targetFunctions, "Open")     // os package
	targetFunctions = append(targetFunctions, "ReadFile") // ioutil package

	// SQL function
	targetFunctions = append(targetFunctions, "Query")    // database/sql package
	targetFunctions = append(targetFunctions, "QueryRow") // database/sql package
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

		for _, targetFunction := range targetFunctions {
			if id != nil && id.Name == targetFunction {
				for _, arg := range args {
					// targetFunction(arg)
					if ident, ok := arg.(*ast.Ident); ok {
						flag = checkTaintNode(ident)
					}
					// targetFunction(arg1 + arg2)
					if binaryExpr, ok := arg.(*ast.BinaryExpr); ok {
						flag = checkTaintNode(binaryExpr)
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

// This function returns true if the argument node is tainted.
// - Search the AST rooted at the argument node.
// - Check if this node is created only from constant values.
func checkTaintNode(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.BasicLit:
		return false
	case *ast.Ident:
		if node.Obj == nil || node.Obj.Kind == ast.Con {
			return false
		}
		if node, ok := node.Obj.Decl.(ast.Node); ok {
			return checkTaintNode(node)
		}
	case *ast.AssignStmt:
		if len(node.Rhs) == 0 {
			return true
		}
		for _, rval := range node.Rhs {
			if checkTaintNode(rval) {
				return true
			}
		}
		return false
	case *ast.BinaryExpr:
		return checkTaintNode(node.X) || checkTaintNode(node.Y)
	case *ast.CallExpr:
		// TODO: We should make sure that the return value of this function is trustworthy,
		//       but for simplicity we don't trust all functions.
		return true
	}
	return false
}
