// Package main provides the CLI for the usebacktick analyzer.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ccoveille/usebacktick/analyzer"
)

func main() {
	var settings any // TODO: add support for settings
	singlechecker.Main(analyzer.New(settings))
}
