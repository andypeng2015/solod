package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// emitStringLit emits a string literal, handling both interpreted and raw strings.
func (g *Generator) emitStringLit(n *ast.BasicLit) {
	fmt.Fprintf(g.state.writer, "so_str(%s)", rawStringValue(n))
}

// rawStringValue returns the C string literal for a Go string literal,
// handling both interpreted and raw strings. Does not include the so_str() wrapper.
func rawStringValue(n *ast.BasicLit) string {
	if strings.HasPrefix(n.Value, "`") {
		// Raw string: strip backticks, escape for C.
		raw := n.Value[1 : len(n.Value)-1]
		var b strings.Builder
		for _, ch := range raw {
			switch ch {
			case '\\':
				b.WriteString(`\\`)
			case '"':
				b.WriteString(`\"`)
			case '\n':
				b.WriteString(`\n`)
			case '\t':
				b.WriteString(`\t`)
			case '\r':
				b.WriteString(`\r`)
			default:
				b.WriteRune(ch)
			}
		}
		return `"` + b.String() + `"`
	}
	return n.Value
}

func stringCompareFunc(op token.Token) string {
	switch op {
	case token.EQL:
		return "so_string_eq"
	case token.NEQ:
		return "so_string_ne"
	case token.LSS:
		return "so_string_lt"
	case token.LEQ:
		return "so_string_lte"
	case token.GTR:
		return "so_string_gt"
	case token.GEQ:
		return "so_string_gte"
	}
	panic("unreachable")
}
