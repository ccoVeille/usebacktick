// Package golangci supports using this analyzer as a golangci-lint plugin.
package golangci

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/ccoveille/usebacktick/analyzer"
)

//nolint:gochecknoinits // init needed for plugin
func init() {
	register.Plugin("usebacktick", func(_ any) (register.LinterPlugin, error) {
		return &usebacktickPlugin{}, nil
	})
}

type usebacktickPlugin struct{}

var _ register.LinterPlugin = new(usebacktickPlugin)

// BuildAnalyzers returns the analyzers to be run by golangci-lint.
//
// This method is part of the [register.LinterPlugin]
func (usebacktickPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}, nil
}

// GetLoadMode returns the load mode for the analyzers to be run by golangci-lint.
//
// This method is part of the [register.LinterPlugin]
func (usebacktickPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
