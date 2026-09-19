package importtest

import "regexp"

type stub struct{}

func (stub) Compile(string) {}

func f() {
	// regexp shadows the regexp package
	regexp := stub{}
	regexp.Compile("a")
}

// keep the import in use without adding a call for isCallTo to match
var _ = regexp.MustCompile
