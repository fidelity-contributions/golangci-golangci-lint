package godox

import (
	"go/token"
	"strings"

	"github.com/matoous/godox"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.GodoxSettings) *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(&analysis.Analyzer{
			Name: "godox",
			Doc:  "Detects usage of FIXME, TODO and other keywords inside comments",
			Run: func(pass *analysis.Pass) (any, error) {
				return run(pass, settings), nil
			},
		}).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(pass *analysis.Pass, settings *config.GodoxSettings) error {
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		messages, err := godox.Run(file, pass.Fset, settings.Keywords...)
		if err != nil {
			return err
		}

		if len(messages) == 0 {
			continue
		}

		nonAdjPosition := pass.Fset.PositionFor(file.Pos(), false)

		ft := pass.Fset.File(file.Pos())

		for _, i := range messages {
			pass.Report(analysis.Diagnostic{
				Pos:     ft.LineStart(goanalysis.AdjustPos(i.Pos.Line, nonAdjPosition.Line, position.Line)) + token.Pos(i.Pos.Column),
				Message: strings.TrimRight(i.Message, "\n"),
			})
		}
	}

	return nil
}
