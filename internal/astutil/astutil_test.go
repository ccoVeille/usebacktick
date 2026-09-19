package astutil

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"slices"
	"testing"
)

func TestFind(t *testing.T) {
	tests := []struct {
		file      string
		path      string
		wantNames []string
	}{
		{file: "plain.go", path: "regexp", wantNames: []string{"regexp"}},
		{file: "subpackage.go", path: "regexp/syntax", wantNames: []string{"syntax"}},
		{file: "alias.go", path: "regexp", wantNames: []string{"re"}},
		{file: "blank.go", path: "regexp", wantNames: nil},
		{file: "dot.go", path: "regexp", wantNames: []string{"."}},
		{file: "missing.go", path: "regexp", wantNames: nil},
		{file: "shadowed.go", path: "regexp", wantNames: []string{"regexp"}},
		{file: "multialias.go", path: "regexp", wantNames: []string{"re", "reg", "regexp"}},
		{file: "aliased_other_package.go", path: "regexp", wantNames: nil},
		{file: "dot_shadowed.go", path: "regexp", wantNames: []string{"."}},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			file := parseTestdata(t, tt.file)

			got := FindImport(file, tt.path)

			gotNames := slices.Clone(got.names)
			slices.Sort(gotNames)

			wantNames := slices.Clone(tt.wantNames)
			slices.Sort(wantNames)

			if !slices.Equal(gotNames, wantNames) {
				t.Fatalf("Find(%q).names = %v, want %v", tt.file, got.names, tt.wantNames)
			}
		})
	}
}

func TestIsCallTo(t *testing.T) {
	tests := []struct {
		file        string
		path        string
		funcNames   []string
		wantMatches int
	}{
		{file: "plain.go", path: "regexp", funcNames: []string{"Compile", "MustCompile"}, wantMatches: 2},
		{file: "subpackage.go", path: "regexp/syntax", funcNames: []string{"Parse"}, wantMatches: 1},
		{file: "subpackage.go", path: "regexp", funcNames: []string{"Parse"}, wantMatches: 0},
		{file: "alias.go", path: "regexp", funcNames: []string{"Compile"}, wantMatches: 1},
		{file: "blank.go", path: "regexp", funcNames: []string{"Compile"}, wantMatches: 0},
		{file: "dot.go", path: "regexp", funcNames: []string{"Compile"}, wantMatches: 1},
		{file: "missing.go", path: "regexp", funcNames: []string{"Compile"}, wantMatches: 0},
		{file: "shadowed.go", path: "regexp", funcNames: []string{"Compile"}, wantMatches: 0},
		{file: "multialias.go", path: "regexp", funcNames: []string{"Compile", "MustCompile"}, wantMatches: 2},
		{file: "aliased_other_package.go", path: "regexp", funcNames: []string{"Compile"}, wantMatches: 0},
		{file: "dot_shadowed.go", path: "regexp", funcNames: []string{"Compile"}, wantMatches: 0},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			file := parseTestdata(t, tt.file)
			imp := FindImport(file, tt.path)

			var matches int

			ast.Inspect(file, func(n ast.Node) bool {
				if imp.IsCallTo(n, tt.funcNames...) {
					matches++
				}

				return true
			})

			if matches != tt.wantMatches {
				t.Fatalf("IsCallTo matches in %q = %d, want %d", tt.file, matches, tt.wantMatches)
			}
		})
	}
}

func parseTestdata(t *testing.T, name string) *ast.File {
	t.Helper()

	path := filepath.Join("testdata", name)

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("ParseFile(%q): %v", path, err)
	}

	return file
}
