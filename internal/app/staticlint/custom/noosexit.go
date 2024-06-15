package custom

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var NoOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "reports the use of os.Exit in main function",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	funcChecker := func(x *ast.FuncDecl) {
		// ищем и проверяем функцию
		for _, stmt := range x.Body.List {
			if callExpr, ok := stmt.(*ast.ExprStmt); ok {
				// проверяем вызов
				if call, ok := callExpr.X.(*ast.CallExpr); ok {
					if s, ok := call.Fun.(*ast.SelectorExpr); ok && s.Sel.Name == "Exit" {
						pass.Reportf(call.Pos(), "avoid using os.Exit in main function")
					}
				}
			}
		}
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				// проверяем пакет
				if x.Name.Name != "main" {
					return false
				}
				for _, obj := range x.Decls {
					// ищем и проверяем функцию
					if x, ok := obj.(*ast.FuncDecl); ok && x.Name.Name == "main" {
						funcChecker(x)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
