package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestUseBackticks(t *testing.T) {
	tests := []struct {
		name     string
		litValue string
		want     string
		wantOK   bool
	}{
		{
			name:     "simple escape",
			litValue: `"foo\"bar"`,
			want:     "`foo\"bar`",
			wantOK:   true,
		},
		{
			name:     "backslash escape",
			litValue: `"foo\\bar"`,
			want:     "`foo\\bar`",
			wantOK:   true,
		},
		{
			name:     "already simplest form",
			litValue: `"foobar"`,
			wantOK:   false,
		},
		{
			name:     "contains tab",
			litValue: `"foo\tbar"`,
			wantOK:   false,
		},
		{
			name:     "not backquotable due to backtick in value",
			litValue: "\"foo`bar\"",
			wantOK:   false,
		},
		{
			name:     "unicode escape sequence preserved",
			litValue: `"foo\u1F601bar"`,
			wantOK:   false,
		},
		{
			name:     "hex escape sequence preserved",
			litValue: `"foo\x41bar"`,
			wantOK:   false,
		},
		{
			name:     "newline not backquotable",
			litValue: `"foo\nbar\"baz"`,
			wantOK:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := useBackticks(tt.litValue)
			if ok != tt.wantOK {
				t.Fatalf("useBackticks(%q) ok = %v, want %v", tt.litValue, ok, tt.wantOK)
			}

			if ok && got != tt.want {
				t.Fatalf("useBackticks(%q) = %q, want %q", tt.litValue, got, tt.want)
			}
		})
	}
}

func TestQuotedStringLiteralValue(t *testing.T) {
	const src = `package p

func f() {
	_ = "double"
	_ = ` + "`raw`" + `
	_ = 42
}
`

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	var lits []*ast.BasicLit

	ast.Inspect(file, func(n ast.Node) bool {
		lit, skip := quotedStringLiteralValue(n)
		if !skip {
			lits = append(lits, lit)
		}

		return true
	})

	if len(lits) != 1 {
		t.Fatalf("got %d double-quoted string literals, want 1", len(lits))
	}

	if lits[0].Value != `"double"` {
		t.Fatalf("got literal %q, want %q", lits[0].Value, `"double"`)
	}
}

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, New(nil), "a")
}
