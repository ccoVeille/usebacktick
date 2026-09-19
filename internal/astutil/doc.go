// Package astutil provides helpers to detect how a package is imported in a file,
// and to recognize calls to specific functions on that package.
//
// Its purpose is to be used by analyzers that need to detect calls to specific functions on a package,
// but without using Go types.
//
// go/types requires loading and type-checking a package's dependencies.
//
// This package only walks the AST (import specs and identifiers), so it works under that load mode,
// at the cost of the heuristics documented on [PackageImport.IsCallTo] (e.g. a shadowing local declaration is enough to hide a match).
//
// # Relation to other tools
//
// [golang.org/x/tools/go/ast/astutil] (same name, different module) solves half of this problem:
// its Imports and UsesImport helpers detect whether/how a package is imported, purely from the AST.
// It has no equivalent of IsCallTo, since it isn't concerned with matching specific function calls.
//
// honnef.co/go/tools (staticcheck), whose code.IsCallToAny matches calls like "regexp.MustCompile"
// the same way IsCallTo does, resolves the call through go/types (pass.TypesInfo), that this package
// avoids for speed and simplicity.
package astutil
