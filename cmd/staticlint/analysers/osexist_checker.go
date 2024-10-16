// Package shorteneranalysers contains custom static analysers for this project.
package shorteneranalysers

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer checks code for "os.Exit()" calls in "main" function of package "main".
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit in func 'main' of package 'main'",
	Run:  osExitAnalyzerRun,
}

// osExitAnalyzerRun is a main func of OsExitAnalyzer witch checks code.
func osExitAnalyzerRun(pass *analysis.Pass) (interface{}, error) {

	var currentPackage string
	var currentFunc string

	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {

			switch nn := n.(type) {
			case *ast.File:
				currentPackage = nn.Name.Name
			case *ast.FuncDecl:
				currentFunc = nn.Name.Name
			case *ast.CallExpr:
				if currentPackage == "main" && currentFunc == "main" {
					if selector, ok := nn.Fun.(*ast.SelectorExpr); ok {
						if pkgIdent, ok := selector.X.(*ast.Ident); ok {
							if pkgIdent.Name == "os" && selector.Sel.Name == "Exit" {
								pass.Reportf(selector.Pos(), "call os.Exit is not allowed in func \"main\" of package \"main\"")
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
