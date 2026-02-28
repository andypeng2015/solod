package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"strings"
)

// emitAssignStmt emits an assignment statement.
func (g *Generator) emitAssignStmt(stmt *ast.AssignStmt) {
	switch stmt.Tok {
	case token.DEFINE:
		w := g.state.writer
		// Detect: _, ok := s.(Rect)
		if len(stmt.Lhs) == 2 && len(stmt.Rhs) == 1 {
			if ta, ok := stmt.Rhs[0].(*ast.TypeAssertExpr); ok {
				g.emitTypeAssertion(w, stmt, ta)
				return
			}
		}
		// Detect multi-return: a, b := vals()
		if len(stmt.Lhs) > 1 && len(stmt.Rhs) == 1 {
			if call, ok := stmt.Rhs[0].(*ast.CallExpr); ok {
				g.emitMultiReturnDefine(w, stmt, call)
				return
			}
		}
		// Regular define: group consecutive variables by type.
		i := 0
		for i < len(stmt.Lhs) {
			ident := stmt.Lhs[i].(*ast.Ident)
			if ident.Name == "_" {
				i++
				continue
			}
			def := g.types.Defs[ident]
			if def == nil {
				// Redeclared variable - emit plain assignment.
				fmt.Fprintf(w, "%s%s = ", g.indent(), ident.Name)
				g.emitExpr(stmt.Rhs[i])
				fmt.Fprintf(w, ";\n")
				i++
				continue
			}
			typ := def.Type()
			cType := g.mapType(stmt, typ)
			fmt.Fprintf(w, "%s%s %s = ", g.indent(), cType, ident.Name)
			g.emitExpr(stmt.Rhs[i])
			i++
			for i < len(stmt.Lhs) {
				nextIdent := stmt.Lhs[i].(*ast.Ident)
				if nextIdent.Name == "_" {
					break
				}
				nextDef := g.types.Defs[nextIdent]
				if nextDef == nil {
					break
				}
				nextCType := g.mapType(stmt, nextDef.Type())
				if nextCType != cType {
					break
				}
				fmt.Fprintf(w, ", %s = ", nextIdent.Name)
				g.emitExpr(stmt.Rhs[i])
				i++
			}
			fmt.Fprintf(w, ";\n")
		}

	case token.ASSIGN:
		w := g.state.writer
		// Detect multi-return: b, a = swap(a, b)
		if len(stmt.Lhs) > 1 && len(stmt.Rhs) == 1 {
			if call, ok := stmt.Rhs[0].(*ast.CallExpr); ok {
				g.emitMultiReturnAssign(w, stmt, call)
				return
			}
		}
		// Regular assignment.
		for i, lhs := range stmt.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok && ident.Name == "_" {
				fmt.Fprintf(w, "%s(void)", g.indent())
				if g.needsVoidParens(stmt.Rhs[i]) {
					fmt.Fprintf(w, "(")
					g.emitExpr(stmt.Rhs[i])
					fmt.Fprintf(w, ")")
				} else {
					g.emitExpr(stmt.Rhs[i])
				}
				fmt.Fprintf(w, ";\n")
				continue
			}
			fmt.Fprintf(w, "%s", g.indent())
			g.emitExpr(lhs)
			fmt.Fprintf(w, " = ")
			g.emitExpr(stmt.Rhs[i])
			fmt.Fprintf(w, ";\n")
		}

	case token.ADD_ASSIGN, token.SUB_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN,
		token.REM_ASSIGN, token.OR_ASSIGN, token.AND_ASSIGN, token.XOR_ASSIGN,
		token.SHL_ASSIGN, token.SHR_ASSIGN:
		w := g.state.writer
		fmt.Fprintf(w, "%s", g.indent())
		g.emitExpr(stmt.Lhs[0])
		fmt.Fprintf(w, " %s ", stmt.Tok)
		g.emitExpr(stmt.Rhs[0])
		fmt.Fprintf(w, ";\n")

	default:
		g.fail(stmt, "unsupported AssignStmt token: %s", stmt.Tok)
	}
}

// emitMultiReturnDefine emits a multi-return define assignment (e.g. a, b := vals()).
// Out-parameters (index 1+) are declared first, then the primary return value
// is declared and initialized with the call result, e.g.:
// `so_Error err; int n = work(a, b, &err);`
func (g *Generator) emitMultiReturnDefine(w io.Writer, stmt *ast.AssignStmt, call *ast.CallExpr) {
	// Declare out-param variables (index 1+), grouped by type.
	type varInfo struct {
		name  string
		cType string
	}

	// Collect out-parameters, skipping blank identifiers and redeclared variables.
	var outVars []varInfo
	for _, lhs := range stmt.Lhs[1:] {
		ident := lhs.(*ast.Ident)
		if ident.Name == "_" {
			continue
		}
		def := g.types.Defs[ident]
		if def == nil {
			continue // redeclared variable
		}
		outVars = append(outVars, varInfo{ident.Name, g.mapType(stmt, def.Type())})
	}

	// Group consecutive out-parameters by type and emit declarations.
	i := 0
	for i < len(outVars) {
		cType := outVars[i].cType
		names := []string{outVars[i].name}
		for i+1 < len(outVars) && outVars[i+1].cType == cType {
			i++
			names = append(names, outVars[i].name)
		}
		fmt.Fprintf(w, "%s%s %s;\n", g.indent(), cType, strings.Join(names, ", "))
		i++
	}

	// Build out-args from LHS vars at index 1+.
	g.state.outArgs = g.emitOutArgs(w, stmt, call)
	defer func() { g.state.outArgs = nil }()

	// Emit the call with first var declaration+initialization.
	firstIdent := stmt.Lhs[0].(*ast.Ident)
	if firstIdent.Name == "_" {
		fmt.Fprintf(w, "%s", g.indent())
		g.emitExpr(call)
		fmt.Fprintf(w, ";\n")
	} else {
		def := g.types.Defs[firstIdent]
		if def != nil {
			cType := g.mapType(stmt, def.Type())
			fmt.Fprintf(w, "%s%s %s = ", g.indent(), cType, firstIdent.Name)
		} else {
			fmt.Fprintf(w, "%s%s = ", g.indent(), firstIdent.Name)
		}
		g.emitExpr(call)
		fmt.Fprintf(w, ";\n")
	}
}

// emitMultiReturnAssign emits a multi-return assignment (e.g. b, a = swap(a, b)).
func (g *Generator) emitMultiReturnAssign(w io.Writer, stmt *ast.AssignStmt, call *ast.CallExpr) {
	// Build out-args from LHS vars at index 1+.
	g.state.outArgs = g.emitOutArgs(w, stmt, call)
	defer func() { g.state.outArgs = nil }()

	// Emit the call.
	firstIdent := stmt.Lhs[0].(*ast.Ident)
	if firstIdent.Name == "_" {
		fmt.Fprintf(w, "%s", g.indent())
		g.emitExpr(call)
		fmt.Fprintf(w, ";\n")
	} else {
		fmt.Fprintf(w, "%s", g.indent())
		g.emitExpr(stmt.Lhs[0])
		fmt.Fprintf(w, " = ")
		g.emitExpr(call)
		fmt.Fprintf(w, ";\n")
	}
}

func (g *Generator) emitOutArgs(w io.Writer, stmt *ast.AssignStmt, call *ast.CallExpr) []string {
	sig := g.types.Types[call.Fun].Type.(*types.Signature)
	var outArgs []string
	for j, lhs := range stmt.Lhs[1:] {
		ident := lhs.(*ast.Ident)
		if ident.Name == "_" {
			g.state.nDiscard++
			name := fmt.Sprintf("_d%d", g.state.nDiscard)
			cType := g.mapType(stmt, sig.Results().At(j+1).Type())
			fmt.Fprintf(w, "%s%s %s;\n", g.indent(), cType, name)
			outArgs = append(outArgs, "&"+name)
			continue
		}
		outArgs = append(outArgs, "&"+ident.Name)
	}
	return outArgs
}
