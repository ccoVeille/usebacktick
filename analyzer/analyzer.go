// Package analyzer implements a static analysis tool that suggests using raw string literals (backticks)
// instead of quoted string literals when it improves readability and does not change the meaning of the string.
package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer defines the linter configuration for the go/analysis framework.
var Analyzer = &analysis.Analyzer{
	Name: "usebacktick",
	URL:  "https://github.com/ccoVeille/usebacktick",
	Doc:  "reports string literals that can be simplified with raw string literals (backticks).",

	Run: func(pass *analysis.Pass) (any, error) {
		for _, file := range pass.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if isRegexpCompileCall(pass, n) {
					// ignore all arguments to the regexp function.
					return false
				}

				lit, skip := quotedStringLiteralValue(n)
				if skip {
					return true
				}

				newLitValue, replaced := useBackticks(lit.Value)
				if !replaced {
					return true
				}

				pass.Report(analysis.Diagnostic{
					Pos:     lit.Pos(),
					Message: "use raw string literal",
					URL:     "https://github.com/ccoVeille/usebacktick",
					SuggestedFixes: []analysis.SuggestedFix{
						{
							Message: "replace with raw string literal",
							TextEdits: []analysis.TextEdit{
								{
									Pos:     lit.Pos(),
									End:     lit.End(),
									NewText: []byte(newLitValue),
								},
							},
						},
					},
				})

				return true
			})
		}
		return nil, nil
	},
}

func isRegexpCompileCall(pass *analysis.Pass, n ast.Node) bool {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if sel.Sel == nil { // This should not happen, unless the language spec changes, but let's be safe.
		return false
	}

	obj := pass.TypesInfo.ObjectOf(sel.Sel)
	if obj == nil {
		return false
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}
	if pkg.Path() != "regexp" {
		return false
	}

	return obj.Name() == "Compile" || obj.Name() == "MustCompile"
}

func quotedStringLiteralValue(n ast.Node) (*ast.BasicLit, bool) {
	lit, ok := n.(*ast.BasicLit)
	if !ok {
		return nil, true
	}

	if lit.Kind != token.STRING {
		return nil, true
	}

	if strings.HasPrefix(lit.Value, "`") {
		// the string literal is already a raw string literal
		return nil, true
	}

	return lit, false
}

// useBackticks checks if a raw string literal (backticks) can and should be suggested for litValue.
// Returns the new string literal and true if backticks are suitable.
func useBackticks(litValue string) (string, bool) {
	unquoted, err := strconv.Unquote(litValue)
	if err != nil {
		return "", false
	}

	if litValue == `"`+unquoted+`"` {
		// the string literal is already in its simplest form
		// using backticks would bring nothing
		return "", false
	}

	// TAB is allowed by [strconv.CanBackquote] but it's not readable in a raw string literal, so we don't want to suggest it
	if strings.Contains(unquoted, "\t") {
		return "", false
	}

	if !strconv.CanBackquote(unquoted) {
		// if the string literal is not convertible to a raw string, it's impossible to use backticks
		return "", false
	}

	for _, r := range []string{`\u`, `\U`, `\x`} {
		if strings.Contains(litValue, r) && !strings.Contains(unquoted, r) {
			// if someone used an escape sequence like \u, \U, \x or \X in the original string literal
			// and we no longer have it in the unquoted string
			// we can assume the user wanted to use the escape sequence and we should not suggest a change to backticks
			//
			// some examples where we can assume the user wanted to use the escape sequence:
			// emoji: "\u1F601" for 😁
			// unicode version of char: "\u0041" for letter a
			// hexadecimal version of char: "\x41" for letter a
			return "", false
		}
	}

	return "`" + unquoted + "`", true
}
