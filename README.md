# usebacktick

`usebacktick` is a Go analyzer that finds double-quoted string literals that can be simplified with raw string literals (backticks).

It reports a diagnostic with an automatic suggested fix.

This linter is not only about style, but also about readability and maintainability.

It is built with Go's `go/analysis` framework and can run standalone or as a `go vet` tool.

Source code: <https://github.com/ccoVeille/usebacktick>

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

- the string contains a tab escape (`\t`):

    Even though raw string literals can represent tabs directly

    ```go
    message = "foo\tbar"   // the tabulation is escaped via the \t
    message = `foo  bar` // the tabulation is not visible
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

- the string contains a newline or backtick:
     
    Raw string literals can contain newlines, but this analyzer does not suggest introducing
    literal newlines because they are typically not readable in this context.
    
    Backticks cannot appear inside raw string literals.
    See: <https://go.dev/ref/spec#String_literals>
    
    For example:

    ```go
    message = "\n"
    message = "`"
    ```

## Installation

```console
go install github.com/ccoveille/usebacktick@latest
```

The project currently requires Go 1.26.4 or newer.

## Usage

When the analyzer reports a finding, editors and analysis tools that support Go suggested fixes can apply the conversion automatically.

The analyzer can also be used by tools that support Go analyzers, including
`go vet`, `golangci-lint`, and compatible editor integrations.

### CLI

Run the analyzer directly against packages:

```console
usebacktick ./...
```

Use `-fix` to fix all issues.

```console
usebacktick -fix ./...
```

### go vet

To run it through `go vet`:

```console
go vet -vettool="$(which usebacktick)" ./...
```

```console
go vet -vettool="$(which usebacktick)" -fix ./...
```

Use `-fix` to fix all issues.

### golangci-lint plugin

`usebacktick` can also be used as a [golangci-lint](https://golangci-lint.run/) module plugin, built into a custom `golangci-lint` binary via [`golangci-lint custom`](https://golangci-lint.run/plugins/module-plugins/).

Add a `.custom-gcl.yml` file listing the plugin:

```yaml
version: v2.13.2
plugins:
  - module: 'github.com/ccoveille/usebacktick'
    import: 'github.com/ccoveille/usebacktick/golangci'
    version: v0.2.0
```

Update `version` (the `golangci-lint` version) and the plugin's `version` to the latest
available releases before building, since `golangci-lint custom` pins exact versions
and does not resolve them automatically.

Then build the custom binary:

```console
golangci-lint custom
```

This produces a `custom-gcl` binary in the current directory, which behaves like `golangci-lint` with `usebacktick` included.

Finally, declare the linter in your `golangci-lint` YAML configuration file:

```yaml
version: "2"

# ...

linters:
  settings:

    # ...

    custom:
      usebacktick:
        type: module
        description: Linter for using backticks instead of double quotes where possible
        original-url: https://github.com/ccoVeille/usebacktick
```

Then use `custom-gcl` as a `golangci-lint` replacement:

```console
./custom-gcl lint --enable-only usebacktick ./...
```

Use `-fix` to fix detected issues.

```console
./custom-gcl lint --enable-only usebacktick -fix ./...
```

## Development

Run the test suite with:

```console
go test ./...
```
