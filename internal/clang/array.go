package clang

import (
	"fmt"
	"go/ast"
	"go/types"
)

// emitArrayLit emits a fixed-size array literal as a so_Slice compound literal.
// Example: [5]int{1, 2, 3, 4, 5} → {(int[5]){1, 2, 3, 4, 5}, 5}
func (g *Generator) emitArrayLit(n *ast.CompositeLit) {
	w := g.state.writer
	arr := g.types.TypeOf(n).Underlying().(*types.Array)
	elemType := g.mapType(n, arr.Elem())
	size := int(arr.Len())
	fmt.Fprintf(w, "(so_Slice){(%s[%d]){", elemType, size)

	if hasKeyedElements(n) {
		g.emitSparseArrayValues(n)
	} else {
		for i, elt := range n.Elts {
			if i > 0 {
				fmt.Fprintf(w, ", ")
			}
			g.emitExpr(elt)
		}
	}

	fmt.Fprintf(w, "}, %d, %d}", size, size)
}

// emitSliceLit emits a slice literal as a so_Slice compound literal.
// Example: []int{1, 2, 3, 4} → {(so_int[4]){1, 2, 3, 4}, 4, 4}
func (g *Generator) emitSliceLit(n *ast.CompositeLit) {
	w := g.state.writer
	sl := g.types.TypeOf(n).Underlying().(*types.Slice)
	elemType := g.mapType(n, sl.Elem())
	size := len(n.Elts)
	fmt.Fprintf(w, "(so_Slice){(%s[%d]){", elemType, size)
	for i, elt := range n.Elts {
		if i > 0 {
			fmt.Fprintf(w, ", ")
		}
		g.emitExpr(elt)
	}
	fmt.Fprintf(w, "}, %d, %d}", size, size)
}

// emitSparseArrayValues emits array values using C99 designated initializers
// for keyed elements. Example: [...]int{100, 3: 400, 500} → 100, [3] = 400, 500
func (g *Generator) emitSparseArrayValues(n *ast.CompositeLit) {
	w := g.state.writer
	for i, elt := range n.Elts {
		if i > 0 {
			fmt.Fprintf(w, ", ")
		}
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			fmt.Fprintf(w, "[")
			g.emitExpr(kv.Key)
			fmt.Fprintf(w, "] = ")
			g.emitExpr(kv.Value)
		} else {
			g.emitExpr(elt)
		}
	}
}

// emitSliceExpr emits a slice expression (e.g. nums[1:4]) as so_slice(nums, int, 1, 4).
func (g *Generator) emitSliceExpr(n *ast.SliceExpr) {
	w := g.state.writer

	// Determine the element type of the slice.
	var elemType string
	switch t := g.types.TypeOf(n.X).Underlying().(type) {
	case *types.Array:
		elemType = g.mapType(n, t.Elem())
	case *types.Slice:
		elemType = g.mapType(n, t.Elem())
	default:
		g.fail(n, "unsupported slice expression type: %T", t)
	}

	// Emit the slice expression as so_slice(x, elemType, low, high).
	fmt.Fprintf(w, "so_slice(")
	g.emitExpr(n.X)
	fmt.Fprintf(w, ", %s, ", elemType)
	if n.Low != nil {
		g.emitExpr(n.Low)
	} else {
		fmt.Fprintf(w, "0")
	}
	fmt.Fprintf(w, ", ")
	if n.High != nil {
		g.emitExpr(n.High)
	} else {
		g.emitExpr(n.X)
		fmt.Fprintf(w, ".len")
	}
	fmt.Fprintf(w, ")")
}

// hasKeyedElements returns true if any element
// in the composite literal uses key:value syntax.
func hasKeyedElements(n *ast.CompositeLit) bool {
	for _, elt := range n.Elts {
		if _, ok := elt.(*ast.KeyValueExpr); ok {
			return true
		}
	}
	return false
}
