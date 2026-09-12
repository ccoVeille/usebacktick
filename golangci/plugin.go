// Package golangci supports using this analyzer as a golangci-lint plugin.
package golangci

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/ccoveille/usebacktick/analyzer"
)

//nolint:gochecknoinits // init needed for plugin
func init() {
	register.Plugin("usebacktick", func(settings any) (register.LinterPlugin, error) {
		return &usebacktickPlugin{
			settings: settings,
		}, nil
	})
}

type usebacktickPlugin struct {
	settings any
}

var _ register.LinterPlugin = new(usebacktickPlugin)

// BuildAnalyzers returns the analyzers to be run by golangci-lint.
//
// This method is part of the [register.LinterPlugin]
func (p usebacktickPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.New(p.settings),
	}, nil
}

// GetLoadMode returns the load mode for the analyzers to be run by golangci-lint.
//
// This method is part of the [register.LinterPlugin]
func (usebacktickPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
