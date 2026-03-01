package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"strings"
)

// emitImports emits #include directives for imports.
func (g *Generator) emitImports(w io.Writer) {
	for _, file := range g.pkg.Syntax {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.IMPORT {
				continue
			}
			for _, spec := range gd.Specs {
				g.emitImportSpec(w, spec.(*ast.ImportSpec))
			}
		}
	}
}

// emitImportSpec emits a #include directive for an import.
func (g *Generator) emitImportSpec(w io.Writer, spec *ast.ImportSpec) {
	path := strings.Trim(spec.Path.Value, `"`)
	// Strip the imported package's own module prefix.
	if imp, ok := g.pkg.Imports[path]; ok && imp.Module != nil {
		path = strings.TrimPrefix(path, imp.Module.Path+"/")
	}
	// Add the package.h file (e.g. package -> package/package.h).
	parts := strings.Split(path, "/")
	parts = append(parts, parts[len(parts)-1]+".h")
	cPath := strings.Join(parts, "/")
	// Emit the #include directive (e.g. #include "package/package.h").
	fmt.Fprintf(w, "#include \"%s\"\n", cPath)
}

// emitHeaderDecls writes declarations for exported package-level symbols.
// Types are emitted first so that const/var and function prototypes
// can reference them.
func (g *Generator) emitHeaderDecls(w io.Writer) {
	// Phase 1: exported types from collected symbols.
	for _, sym := range g.symbols {
		if !sym.exported || sym.kind != symbolType {
			continue
		}
		g.emitTypeSpec(w, sym.typeSpec)
	}
	// Phase 2: const/var declarations from the AST.
	for _, file := range g.pkg.Syntax {
		for _, decl := range file.Decls {
			if gd, ok := decl.(*ast.GenDecl); ok {
				g.emitHeaderGenDecl(w, gd)
			}
		}
	}
	// Phase 3: exported function/method prototypes from collected symbols.
	for _, sym := range g.symbols {
		if !sym.exported || sym.kind == symbolType {
			continue
		}
		fn := newFuncDecl(g, sym.funcDecl)
		fmt.Fprintf(w, "%s %s(%s);\n", fn.returnType(), fn.name(), fn.params())
	}
	fmt.Fprintf(w, "\n")
}

// emitHeaderGenDecl emits extern const/var declarations.
// Type declarations are handled separately via collected symbols.
func (g *Generator) emitHeaderGenDecl(w io.Writer, decl *ast.GenDecl) {
	if hasExternDirective(decl.Doc) {
		return
	}
	if decl.Tok == token.TYPE {
		// Types are handled separately in [Generator.emitHeaderDecls].
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
