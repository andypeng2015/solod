package clang

import (
	"go/ast"
	"go/types"
	"strings"
)

// collectFileExterns collects extern symbols from a single file's declarations.
func (g *Generator) collectFileExterns(pkgName string, file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if !hasExternDirective(d.Doc) {
				continue
			}
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					g.markExtern(pkgName, s.Name.Name)
				case *ast.ValueSpec:
					for _, name := range s.Names {
						g.markExtern(pkgName, name.Name)
					}
				}
			}
		case *ast.FuncDecl:
			if d.Body == nil || hasExternDirective(d.Doc) {
				g.markExtern(pkgName, externFuncKey(d))
			}
		}
	}
}

// isExternCall reports whether a call expression targets an extern C function.
func (g *Generator) isExternCall(call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		// Local package call.
		return g.hasExtern("", fun.Name)
	case *ast.SelectorExpr:
		// Package-qualified call (e.g. stdio.Printf).
		if ident, ok := fun.X.(*ast.Ident); ok {
			if pkgName, ok := g.types.Uses[ident].(*types.PkgName); ok {
				return g.hasExtern(pkgName.Name(), fun.Sel.Name)
			}
		}
	}
	return false
}

// markExtern marks a symbol in a package as extern.
func (g *Generator) markExtern(pkgName, name string) {
	if pkgName != "" {
		name = pkgName + "." + name
	}
	g.externs[name] = true
}

// hasExtern reports whether a symbol in a package is marked as extern.
func (g *Generator) hasExtern(pkgName, name string) bool {
	if pkgName != "" {
		name = pkgName + "." + name
	}
	return g.externs[name]
}

// externFuncKey returns a map key for a function or method declaration.
// Functions use their bare name (e.g. "Foo"), while methods use
// "ReceiverType.Name" (e.g. "T.Foo") to avoid collisions.
func externFuncKey(decl *ast.FuncDecl) string {
	if decl.Recv != nil {
		return recvTypeName(decl.Recv.List[0]) + "." + decl.Name.Name
	}
	return decl.Name.Name
}

// hasExternDirective checks if a comment group contains the //so:extern directive.
func hasExternDirective(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	for _, c := range doc.List {
		if strings.TrimSpace(c.Text) == "//so:extern" {
			return true
		}
	}
	return false
}
