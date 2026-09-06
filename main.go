// Package main provides the CLI for the usebacktick analyzer.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ccoveille/usebacktick/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
