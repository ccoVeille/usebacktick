package importtest

import "regexp/syntax"

func f() {
	_, _ = syntax.Parse("a", 0)
}