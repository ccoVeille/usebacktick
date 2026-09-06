package a

import (
	"encoding/json"
	"regexp"
	"strings"
)

func target() {
	// strings that could be simplified

	_ = "foo\"bar" // want `raw string literal could improve readability`
	_ = "foo\\bar" // want `raw string literal could improve readability`

	_ = strings.Contains(
		"ab",
		"\\n", // want `raw string literal could improve readability`
	)

	var dest any
	_ = json.Unmarshal(
		[]byte("{foo: \"bar\"}"), // want `raw string literal could improve readability`
		&dest,
	)
}

func stable() {
	// using backticks would be useless as the string literal is already in its simplest form
	_ = "foo bar"

	// already using backticks
	_ = `foo bar`

	// cannot use backticks as the string contains a backtick
	_ = "foo`bar"

	// using backticks is not possible as the string contains new line characters
	_ = "foo\nbar"
}

func notReadable() {
	// even if tab can be represented in a raw string literal, it's not readable, so we don't want to suggest it
	_ = "foo\tbar"

	// Unicode and hex escape sequences are preserved
	// if someone wants to use them, we don't want to suggest to use backticks
	_ = "foo\u1F601bar"
	_ = "foo\x41bar"
	_ = "foo\u1F601bar"
}

// regexpCompile is already reported with staticcheck
// https://staticcheck.dev/docs/checks/#S1007
func regexpCompile() {
	_, _ = regexp.Compile("foo\\d+bar")
	_ = regexp.MustCompile("foo\\d+bar")
}
