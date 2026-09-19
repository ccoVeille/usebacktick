package astutil

import (
	"go/ast"
	"path"
	"slices"
	"strconv"
)

// dotQualifier is used as a stand-in name for a dot-imported package, whose functions are called unqualified.
const dotQualifier = "."

// PackageImport describes how a package is imported in a file, and the local names used to reference it.
// A file may import the same package path more than once under different aliases.
// A dot import is represented by dotQualifier, since its functions are called unqualified.
type PackageImport struct {
	names []string
}

// FindImport looks for path among file's imports and returns the local names it's referenced by.
func FindImport(file *ast.File, packageName string) PackageImport {
	var p PackageImport

	for _, imp := range file.Imports {
		importPath, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}

		if importPath != packageName {
			continue
		}

		if imp.Name == nil {
			p.names = append(p.names, path.Base(importPath))
			continue
		}

		switch imp.Name.Name {
		case "_":
			// ignore blank imports
		case ".":
			p.names = append(p.names, dotQualifier)
		default:
			// use the alias name for the package
			p.names = append(p.names, imp.Name.Name)
		}
	}

	return p
}

// IsCallTo reports whether n is a call to one of funcNames on the imported package,
// under any of its local names, or unqualified if the package is dot-imported.
func (p PackageImport) IsCallTo(n ast.Node, funcNames ...string) bool {
	if len(p.names) == 0 {
		return false
	}

	call, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}

	qualifier, name, ok := callTarget(call)
	if !ok {
		return false
	}

	if !slices.Contains(p.names, qualifier) {
		return false
	}

	return slices.Contains(funcNames, name)
}

// callTarget extracts the package qualifier (dotQualifier for an unqualified call) and function name from call.
func callTarget(call *ast.CallExpr) (qualifier, name string, ok bool) {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		if fun.Obj != nil {
			// this is not the dot-imported function, it could be a local declaration with the same name
			return "", "", false
		}

		return dotQualifier, fun.Name, true
	case *ast.SelectorExpr:
		if fun.Sel == nil { // This should not happen, unless the language spec changes, but let's be safe.
			return "", "", false
		}

		pkg, ok := fun.X.(*ast.Ident)
		if !ok {
			return "", "", false
		}

		if pkg.Obj != nil {
			// this is not a package, it could be a variable with the same name as the package
			return "", "", false
		}

		return pkg.Name, fun.Sel.Name, true
	default:
		return "", "", false
	}
}
