package gochecknoinits

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(&analysis.Analyzer{
			Name:     "gochecknoinits",
			Doc:      "Checks that no init functions are present in Go code",
			Run:      run,
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		}).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(decl ast.Node) {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			return
		}

		fnName := funcDecl.Name.Name
		if fnName == "init" && funcDecl.Recv.NumFields() == 0 {
			pass.Reportf(funcDecl.Pos(), "don't use %s function", internal.FormatCode(fnName))
		}
	})

	return nil, nil
}
