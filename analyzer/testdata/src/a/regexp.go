package a

import (
	"io"
	"regexp"
	"regexp/syntax"
)

// What is already reported by staticcheck S1007
// https://staticcheck.dev/docs/checks/#S1007
func regexpStaticCheckS1007() {
	_, _ = regexp.Compile("foo\\d+bar")
	_ = regexp.MustCompile("foo\\d+bar")
}

// What is not reported by staticcheck S1007
// but that might be reported by staticcheck at some point in the future
// by extending the S1007 check to other regexp functions
// right now, usebacktick reports them.
func regexpStaticCheckS1007NotSupported() {
	_, _ = regexp.CompilePOSIX("foo\\d+bar")  // want `use raw string literal`
	_ = regexp.MustCompilePOSIX("foo\\d+bar") // want `use raw string literal`

	_, _ = syntax.Parse("foo\\d+bar", syntax.POSIX) // want `use raw string literal`

	_, _ = regexp.Match("foo\\d+bar", []byte("foo123bar")) // want `use raw string literal`
	_, _ = regexp.MatchString("foo\\d+bar", "foo123bar")   // want `use raw string literal`

	_ = regexp.QuoteMeta("foo \\n") // want `use raw string literal`

	var rr io.RuneReader
	_, _ = regexp.MatchReader("foo\\d+bar", rr) // want `use raw string literal`

	var f regexp.Regexp
	f.UnmarshalText([]byte("foo\\d+bar")) // want `use raw string literal`
}
