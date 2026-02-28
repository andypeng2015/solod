package clang

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// EmitOptions holds the options for code generation.
type EmitOptions struct {
	Pkg    *packages.Package
	OutDir string
}

// Emit generates C code for the given Go package and all its subpackages,
// and writes it to the specified output directory. Creates a single header
// file with typedefs (.h) and a single implementation file (.c) for each package.
func Emit(opts EmitOptions) error {
	if err := os.MkdirAll(opts.OutDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %s: %w", opts.OutDir, err)
	}
	g := newGenerator(opts.Pkg)
	g.collectExterns()
	if err := g.emitHeader(opts.OutDir); err != nil {
		return err
	}
	return g.emitImpl(opts.OutDir)
}

// State holds the code generation state for the current scope.
type State struct {
	writer io.Writer

	// Current indentation level (number of tabs).
	indent int
	// Current receiver name (for -> access in methods).
	recvName string
	// Out-parameter names for current function (used by emitReturnStmt).
	outParams []string
	// Out-parameter args to append to current call
	// (used by emitFuncCall/emitMethodCall).
	outArgs  []string
	nDiscard int // counter for discard variable names (_d1, _d2, ...)
}

func (s *State) enterFunc(decl FuncDecl) {
	s.outParams = decl.outParams()
	s.nDiscard = 0
}

func (s *State) exitFunc() {
	s.outParams = nil
	s.nDiscard = 0
}

// Generator is responsible for generating C code from Go ASTs.
type Generator struct {
	pkg      *packages.Package
	types    *types.Info
	state    State
	externs  map[string]bool // symbols provided by C headers
	includes []string        // #include directives from comments
	panicked bool            // true after first panic caught in Visit
}

// newGenerator creates a new Generator instance.
func newGenerator(pkg *packages.Package) *Generator {
	return &Generator{
		pkg:     pkg,
		types:   pkg.TypesInfo,
		externs: make(map[string]bool),
	}
}

// emitHeader creates the .h file with typedefs, includes, and extern declarations.
func (g *Generator) emitHeader(dir string) error {
	hName := g.pkg.Name + ".h"
	hPath := filepath.Join(dir, hName)
	hFile, err := os.Create(hPath)
	if err != nil {
		return fmt.Errorf("create header file %s: %w", hPath, err)
	}
	defer hFile.Close()
	fmt.Fprintf(hFile, "#include \"so.h\"\n")
	g.emitHeaderDecls(hFile)
	return nil
}

// emitImpl creates the .c implementation file by walking the AST.
func (g *Generator) emitImpl(dir string) error {
	cName := g.pkg.Name + ".c"
	cPath := filepath.Join(dir, cName)
	cFile, err := os.Create(cPath)
	if err != nil {
		return fmt.Errorf("create C file %s: %w", cPath, err)
	}
	defer cFile.Close()
	fmt.Fprintf(cFile, "#include \"%s.h\"\n", g.pkg.Name)
	// Emit additional #include directives collected from comments.
	for _, inc := range g.includes {
		fmt.Fprintf(cFile, "%s\n", inc)
	}
	g.state.writer = cFile
	for _, file := range g.pkg.Syntax {
		ast.Walk(g, file)
	}
	return nil
}

// collectExterns scans all files for extern symbols and #include directives.
// Body-less functions and declarations annotated with //so:extern are treated
// as external C symbols that should not be emitted.
func (g *Generator) collectExterns() {
	for _, file := range g.pkg.Syntax {
		// Collect // #include comments from the file.
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				text := strings.TrimPrefix(c.Text, "// ")
				if strings.HasPrefix(text, "#include") {
					g.includes = append(g.includes, text)
				}
			}
		}

		// Collect extern symbols from declarations.
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if !hasExternDirective(d.Doc) {
					continue
				}
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						g.externs[s.Name.Name] = true
					case *ast.ValueSpec:
						for _, name := range s.Names {
							g.externs[name.Name] = true
						}
					}
				}
			case *ast.FuncDecl:
				if d.Body == nil {
					g.externs[d.Name.Name] = true
				}
			}
		}
	}
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

// indent returns the current indentation string based on the indent level.
func (g *Generator) indent() string {
	return strings.Repeat("    ", g.state.indent)
}
