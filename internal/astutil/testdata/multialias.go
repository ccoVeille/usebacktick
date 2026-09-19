package importtest

import (
	"regexp"
	re "regexp"
	reg "regexp"
)

func f() {
	_, _ = re.Compile("a")
	_ = reg.MustCompile("a")
	_ = regexp.QuoteMeta("a")
}
