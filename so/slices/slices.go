// Package slices provides various functions useful with slices of any type.
package slices

import (
	_ "embed"

	"github.com/nalgeon/solod/so/mem"
)

//so:embed slices.h
var slices_h string

// Append appends elements to a heap-allocated slice, growing it if needed.
// Returns the updated slice or panics on allocation failure.
// If the allocator is nil, uses the system allocator.
//
//so:extern
func Append[T any](a mem.Allocator, s []T, elems ...T) []T {
	return append(s, elems...)
}

// Extend appends all elements from another heap-allocated slice, growing if needed.
// Returns the updated slice or panics on allocation failure.
// If the allocator is nil, uses the system allocator.
//
//so:extern
func Extend[T any](a mem.Allocator, s []T, other []T) []T {
	return append(s, other...)
}
