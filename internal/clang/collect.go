package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"strings"
)

type symbolKind int

const (
	symbolFunc symbolKind = iota
	symbolMethod
	symbolType
)

type symbol struct {
	kind     symbolKind
	exported bool
	genDecl  *ast.GenDecl // parent GenDecl (for type symbols, enables comment lookup)
	typeSpec *ast.TypeSpec
	funcDecl *ast.FuncDecl
}

// collectSymbols gathers all top-level type and function declarations
// into an ordered list. This list drives both header emission (exported
// symbols) and forward declarations in the .c file (unexported symbols).
func (g *Generator) collectSymbols() {
	// Collect top-level types and functions and their export status.
	for _, file := range g.pkg.Syntax {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				if hasExternDirective(d.Doc) {
					continue
				}
				for _, spec := range d.Specs {
					ts := spec.(*ast.TypeSpec)
					if g.hasExtern("", ts.Name.Name) {
						continue
					}
					g.symbols = append(g.symbols, symbol{
						kind:     symbolType,
						exported: ast.IsExported(ts.Name.Name),
						genDecl:  d,
						typeSpec: ts,
					})
				}
			case *ast.FuncDecl:
				if d.Body == nil || d.Name.Name == "main" {
					continue
				}
				if g.hasExtern("", externFuncKey(d)) {
					continue
				}
				kind := symbolFunc
				exported := ast.IsExported(d.Name.Name)
				if d.Recv != nil {
					kind = symbolMethod
					if exported {
						exported = ast.IsExported(recvTypeName(d.Recv.List[0]))
					}
				}
				g.symbols = append(g.symbols, symbol{
					kind:     kind,
					exported: exported,
					funcDecl: d,
				})
			}
		}
	}

	// Validate that exported functions don't use unexported types.
	for _, sym := range g.symbols {
		if !sym.exported || sym.kind == symbolType {
			continue
		}
		decl := sym.funcDecl
		if g.hasUnexportedTypes(decl) {
			g.fail(decl.Name, "exported function %s uses unexported types", decl.Name.Name)
		}
	}
}

// collectExterns scans all files for extern symbols and include directives.
// Body-less functions and declarations annotated with //so:extern are treated
// as external C symbols that should not be emitted.
func (g *Generator) collectExterns() {
	for _, file := range g.pkg.Syntax {
		// Collect include directives from the file.
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				if include, ok := strings.CutPrefix(c.Text, "//so:include"); ok {
					g.includes = append(g.includes, strings.TrimSpace(include))
				}
			}
		}

		// Collect extern symbols from declarations.
		g.collectFileExterns("", file)
	}

	// Collect externs from imported packages so that isExternCall
	// can identify cross-package extern calls (e.g. stdio.Printf).
	for _, imp := range g.pkg.Imports {
		for _, file := range imp.Syntax {
			g.collectFileExterns(imp.Name, file)
		}
	}
}

// emitForwardDecls writes forward declarations for all unexported symbols.
// Types are emitted first, then functions/methods, so that type names
// are known before function prototypes reference them.
func (g *Generator) emitForwardDecls(w io.Writer) {
	// First pass: unexported types.
	var specs []*ast.TypeSpec
	for _, sym := range g.symbols {
		if sym.exported || sym.kind != symbolType {
			continue
		}
		specs = append(specs, sym.typeSpec)
	}
	if len(specs) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "// -- Forward declarations (types) --")
		for _, spec := range specs {
			g.emitForwardTypeDecl(w, spec)
		}
	}

	// Second pass: unexported functions/methods.
	var funcDecls []*ast.FuncDecl
	for _, sym := range g.symbols {
		if sym.exported || sym.kind == symbolType {
			continue
		}
		funcDecls = append(funcDecls, sym.funcDecl)
	}
	if len(funcDecls) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "// -- Forward declarations (functions and methods) --")
		for _, decl := range funcDecls {
			g.emitFuncProto(w, decl)
			fmt.Fprintln(w, ";")
		}
	}
}

// emitForwardTypeDecl writes a forward declaration for a type.
func (g *Generator) emitForwardTypeDecl(w io.Writer, spec *ast.TypeSpec) {
	cName := g.symbolName(spec.Name.Name)
	switch spec.Type.(type) {
	case *ast.StructType:
		fmt.Fprintf(w, "typedef struct %s %s;\n", cName, cName)
	case *ast.InterfaceType:
		iface := g.types.Defs[spec.Name].Type().Underlying().(*types.Interface)
		if iface.Empty() {
			cType := g.mapType(spec, iface)
			fmt.Fprintf(w, "typedef %s %s;\n", cType, cName)
		} else {
			fmt.Fprintf(w, "typedef struct %s %s;\n", cName, cName)
		}
	case *ast.FuncType:
		named := g.types.Defs[spec.Name].Type().(*types.Named)
		sig := named.Underlying().(*types.Signature)
		retType := g.returnType(spec, sig)
		var params []string
		for p := range sig.Params().Variables() {
			params = append(params, g.mapType(spec, p.Type()))
		}
		fmt.Fprintf(w, "typedef %s (*%s)(%s);\n", retType, cName, strings.Join(params, ", "))
	default:
		typ := g.types.Defs[spec.Name].Type()
		ct := g.mapCType(spec, typ.Underlying())
		fmt.Fprintf(w, "typedef %s;\n", ct.Decl(cName))
	}
}
