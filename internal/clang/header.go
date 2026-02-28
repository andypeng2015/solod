package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
)

// emitHeaderDecls writes declarations for exported package-level symbols.
func (g *Generator) emitHeaderDecls(w io.Writer) {
	for _, file := range g.pkg.Syntax {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				g.emitHeaderGenDecl(w, d)
			case *ast.FuncDecl:
				g.emitHeaderFuncDecl(w, d)
			}
		}
	}
	fmt.Fprintf(w, "\n")
}

// emitHeaderGenDecl emits extern const/var declarations and struct typedefs.
func (g *Generator) emitHeaderGenDecl(w io.Writer, decl *ast.GenDecl) {
	if hasExternDirective(decl.Doc) {
		return
	}
	if decl.Tok == token.TYPE {
		// Exported type declarations (e.g. struct definitions).
		for _, spec := range decl.Specs {
			ts := spec.(*ast.TypeSpec)
			if !ast.IsExported(ts.Name.Name) {
				continue // unexported types are emitted in the .c file
			}
			g.emitTypeSpec(w, ts)
		}
		return
	}

	// Variable and constant declarations.
	for _, spec := range decl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, name := range vs.Names {
			if !ast.IsExported(name.Name) {
				continue
			}
			typ := g.types.Defs[name].Type()
			cType := g.mapType(spec, typ)
			cName := g.symbolName(name.Name)
			switch decl.Tok {
			case token.CONST:
				fmt.Fprintf(w, "extern const %s %s;\n", cType, cName)
			case token.VAR:
				fmt.Fprintf(w, "extern %s %s;\n", cType, cName)
			}
		}
	}
}

// emitHeaderFuncDecl emits a function prototype for exported functions and methods.
func (g *Generator) emitHeaderFuncDecl(w io.Writer, decl *ast.FuncDecl) {
	if decl.Body == nil {
		// Functions with no body are considered externs and ignored.
		return
	}
	if !ast.IsExported(decl.Name.Name) {
		return
	}
	// Methods on unexported types are static, not exported in the header.
	if decl.Recv != nil && !ast.IsExported(recvTypeName(decl.Recv.List[0])) {
		return
	}
	fn := newFuncDecl(g, decl)
	fmt.Fprintf(w, "%s %s(%s);\n", fn.returnType(), fn.name(), fn.params())
}
