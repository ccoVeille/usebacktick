// Package main demonstrates the use of the usebacktick linter.
package main

import (
	"fmt"
	"os"
)

func main() {
	const jsonData = "{ \"name\": \"John Doe\", \"john.doe@example.com\" }"

	_, _ = fmt.Fprintln(os.Stderr, jsonData)
}
