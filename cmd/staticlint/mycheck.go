// multichecker статических анализаторов.
//
// Анализаторы из модуля analysis/passes.
//
//   - defers: Package defers defines an Analyzer that checks for common mistakes in defer statements.
//   - errorsas: The errorsas package defines an Analyzer that checks that the second argument to errors.As is a pointer to a type implementing error.
//   - printf: Package printf defines an Analyzer that checks consistency of Printf format strings and arguments.
//   - shadow: Package shadow defines an Analyzer that checks for shadowed variables.
//   - structtag: Package structtag defines an Analyzer that checks struct field tags are well formed.
//
// Анализаторы из модуля staticcheck.
//   - all SA
//   - all S1
//   - all ST1 -ST1000, -ST1003, -ST1016, -ST1020, -ST1021, -ST1022, -ST1023
//
// Third party анализаторы
//   - github.com/go-critic/go-critic
//   - github.com/julz/importas
//
// Cобственный анализатор
//   - ExitCheckAnalyzer: запрещающает использовать прямой вызов os.Exit в функции main пакета main.
//
// Для запуска выполнить команды
//   - go build cmd/staticlint/mycheck.go
//   - ./mycheck ./...
package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"

	critic "github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/julz/importas"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// ExitCheckAnalyzer - анализатор, запрещающий вызов os.Exit в функции main пакета main
var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for os.Exit in main functions",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// проверяем, какой конкретный тип лежит в узле
			if file.Name.Name != "main" {
				return false
			}
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return false
				}

			case *ast.CallExpr:
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					if s.Sel.Name == "Exit" {
						pass.Reportf(x.Pos(), "Exit in main function")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func main() {
	mychecks := []*analysis.Analyzer{
		ExitCheckAnalyzer,
		critic.Analyzer,
		importas.Analyzer,
		defers.Analyzer,
		errorsas.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	for _, v := range simple.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	styleChecksExclude := map[string]bool{
		"ST1000": false, "ST1003": false, "ST1016": false,
		"ST1020": false, "ST1021": false, "ST1022": false,
		"ST1023": false,
	}
	for _, v := range stylecheck.Analyzers {
		_, ok := styleChecksExclude[v.Analyzer.Name]
		if !ok {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
}
