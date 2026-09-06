// Package main provides the CLI entry point for the usebacktick analyzer.
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ccoveille/usebacktick/internal/analyzer"
)

// Analyzer defines the linter configuration for the go/analysis framework.
var Analyzer = &analysis.Analyzer{
	Name: "backtickliteral",
	Doc:  "détecte les chaînes en double-quotes contenant des échappements qui gagneraient à utiliser des backticks",
	Run:  analyzer.Run,
}

func main() {
	// Runs the linter directly from the command line: go run checker.go ./...
	singlechecker.Main(Analyzer)
}
