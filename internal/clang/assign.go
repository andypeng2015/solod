package clang

import (
	"fmt"
	"go/ast"
	"go/token"
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
			typ := g.types.Defs[ident].Type()
			cType := g.mapType(stmt, typ)
			fmt.Fprintf(w, "%s%s %s = ", g.indent(), cType, ident.Name)
			g.emitExpr(stmt.Rhs[i])
			i++
			for i < len(stmt.Lhs) {
				nextIdent := stmt.Lhs[i].(*ast.Ident)
				nextCType := g.mapType(stmt, g.types.Defs[nextIdent].Type())
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
func (g *Generator) emitMultiReturnDefine(w io.Writer, stmt *ast.AssignStmt, call *ast.CallExpr) {
	// Emit declarations for non-blank LHS vars, grouped by type.
	type varInfo struct {
		name  string
		cType string
	}
	var vars []varInfo
	for _, lhs := range stmt.Lhs {
		ident := lhs.(*ast.Ident)
		if ident.Name == "_" {
			continue
		}
		typ := g.types.Defs[ident].Type()
		vars = append(vars, varInfo{ident.Name, g.mapType(stmt, typ)})
	}
	i := 0
	for i < len(vars) {
		cType := vars[i].cType
		names := []string{vars[i].name}
		for i+1 < len(vars) && vars[i+1].cType == cType {
			i++
			names = append(names, vars[i].name)
		}
		fmt.Fprintf(w, "%s%s %s;\n", g.indent(), cType, strings.Join(names, ", "))
		i++
	}

	// Build out-args from LHS vars at index 1+.
	var outArgs []string
	for _, lhs := range stmt.Lhs[1:] {
		ident := lhs.(*ast.Ident)
		outArgs = append(outArgs, "&"+ident.Name)
	}

	// Emit the call.
	g.state.outArgs = outArgs
	defer func() { g.state.outArgs = nil }()
	firstIdent := stmt.Lhs[0].(*ast.Ident)
	if firstIdent.Name == "_" {
		fmt.Fprintf(w, "%s", g.indent())
		g.emitExpr(call)
		fmt.Fprintf(w, ";\n")
	} else {
		fmt.Fprintf(w, "%s%s = ", g.indent(), firstIdent.Name)
		g.emitExpr(call)
		fmt.Fprintf(w, ";\n")
	}
}

// emitMultiReturnAssign emits a multi-return assignment (e.g. b, a = swap(a, b)).
func (g *Generator) emitMultiReturnAssign(w io.Writer, stmt *ast.AssignStmt, call *ast.CallExpr) {
	// Build out-args from LHS vars at index 1+.
	var outArgs []string
	for _, lhs := range stmt.Lhs[1:] {
		ident := lhs.(*ast.Ident)
		outArgs = append(outArgs, "&"+ident.Name)
	}

	// Emit the call.
	g.state.outArgs = outArgs
	defer func() { g.state.outArgs = nil }()
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
