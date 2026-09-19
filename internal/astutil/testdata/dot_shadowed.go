package importtest

import . "regexp"

func f() {
	// Compile shadows the dot-imported Compile
	Compile := func(string) {}
	Compile("a")
}
