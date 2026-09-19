package importtest

import "regexp"

func f() {
	_, _ = regexp.Compile("a")
	_ = regexp.MustCompile("a")
	_ = regexp.QuoteMeta("a")
}
