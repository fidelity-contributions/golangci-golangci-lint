package processors

import (
	"fmt"
	"path/filepath"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

var _ Processor = (*GeneratedFileFilter)(nil)

type fileSummary struct {
	generated bool
}

// GeneratedFileFilter filters generated files.
type GeneratedFileFilter struct {
	debugf logutils.DebugFunc

	mode    string
	matcher *GeneratedFileMatcher

	fileSummaryCache map[string]*fileSummary
}

func NewGeneratedFileFilter(mode string) *GeneratedFileFilter {
	return &GeneratedFileFilter{
		debugf: logutils.Debug(logutils.DebugKeyGeneratedFileFilter),

		mode:    mode,
		matcher: NewGeneratedFileMatcher(mode),

		fileSummaryCache: map[string]*fileSummary{},
	}
}

func (*GeneratedFileFilter) Name() string {
	return "generated_file_filter"
}

func (p *GeneratedFileFilter) Process(issues []*result.Issue) ([]*result.Issue, error) {
	if p.mode == config.GeneratedModeDisable {
		return issues, nil
	}

	return filterIssuesErr(issues, p.shouldPassIssue)
}

func (*GeneratedFileFilter) Finish() {}

func (p *GeneratedFileFilter) shouldPassIssue(issue *result.Issue) (bool, error) {
	if filepath.Base(issue.FilePath()) == "go.mod" {
		return true, nil
	}

	// The file is already known.
	fs := p.fileSummaryCache[issue.FilePath()]
	if fs != nil {
		return !fs.generated, nil
	}

	fs = &fileSummary{}
	p.fileSummaryCache[issue.FilePath()] = fs

	var err error
	fs.generated, err = p.matcher.IsGeneratedFile(issue.FilePath(), nil)
	if err != nil {
		return false, fmt.Errorf("failed to get doc (%s) of file %s: %w", p.mode, issue.FilePath(), err)
	}

	p.debugf("file %q is generated: %t", issue.FilePath(), fs.generated)

	// don't report issues for autogenerated files
	return !fs.generated, nil
}
