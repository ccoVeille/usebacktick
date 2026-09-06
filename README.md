# usebacktick

`usebacktick` is a Go analyzer that finds double-quoted string literals that can be simplified with raw string literals (backticks).

It reports a diagnostic with an automatic suggested fix.

This linter is not only about style, but also about readability and maintainability.

It is built with Go's `go/analysis` framework and can run standalone or as a `go vet` tool.

Source code: <https://github.com/ccoveille/usebacktick>

## Example

```go
message := "foo\"bar"
if strings.Contains(message, "\\n") {
    // whatever
}
```

becomes:

```go
message := `foo"bar`
if strings.Contains(message, `\n`) {
    // whatever
}
```

## What it checks

The analyzer suggests backticks when all of the following are true:

- the literal is double-quoted;
- converting its value to a raw string is valid; and
- the conversion removes unnecessary escaping.

It does not suggest conversions when:

- the string contains a newline or backtick.

     Raw string literals can contain newlines, but this analyzer does not suggest
     them because control characters are not considered readable in this context.
     Backticks cannot appear inside raw string literals.

    ```go
     message = "\n"
     message = "`"
    ```

- the string contains a tab escape (`\t`):

    Even though raw string literals can represent tabs directly

    ```go
    message = "foo\tbar"   // the tabulation is escaped via the \t
    message = `foo    bar` // the tabulation is not visible
    ```

- the original literal uses an explicit Unicode representation:

    Example:
    
    ```go
    message = "\u1F60 bar" // this is the greek letter omega in lowercase followed by bar
    ```

    Here, we assume the user wanted to use this representation.
    It is out of scope for this linter to suggest that users could have used

    ```go
    message = `ὠ bar`
    ```

- the original literal uses an explicit hexadecimal representation:

    Example:
    
    ```go
    message = "\x41 bar" // this is the A but written in hexa decimal
    ```

    Here, we assume the user wanted to use this representation.
    It is out of scope for this linter to suggest that users could have used

    ```go
    message = `A bar`
    ```

- the string is an argument to `regexp.Compile` or `regexp.MustCompile`

    The use of raw string literals for these methods is already reported by
    Staticcheck's `S1007` check, so this analyzer skips them to avoid duplicate
    diagnostics.

    The `regexp` package accepts these notations as equivalent:
    ```go
    _ = regexp.MustCompile("\t")
    _ = regexp.MustCompile(`\t`)
    ```

## Installation

```console
go install github.com/ccoveille/usebacktick@latest
```

The project currently requires Go 1.26.4 or newer.

## Usage

Run the analyzer directly against packages:

```console
usebacktick ./...
```

To run it through `go vet`:

```console
go vet -vettool="$(which usebacktick)" ./...
```

When the analyzer reports a finding, editors and analysis tools that support Go suggested fixes can apply the conversion automatically.

The analyzer can also be used by tools that support Go analyzers, including
`go vet` and compatible editor integrations.

## Development

Run the test suite with:

```console
go test ./...
```
