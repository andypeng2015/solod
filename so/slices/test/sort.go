package main

import (
	"solod.dev/so/cmp"
	"solod.dev/so/mem"
	"solod.dev/so/slices"
	"solod.dev/so/testing"
)

func descInt(a, b any) int {
	va := *a.(*int)
	vb := *b.(*int)
	return vb - va
}

var sortInts = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
var sortFloat64s = [...]float64{74.3, 59.0, 238.2, -784.0, 2.3, 9845.768, -959.7485, 905, 7.8, 7.8, 74.3, 59.0, 238.2, -784.0, 2.3}
var sortStrs = [...]string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***"}

func TestIsSorted(t *testing.T) {
	// IsSorted: false on unsorted data.
	if slices.IsSorted(sortInts[:]) {
		t.Error("IsSorted: unsorted ints")
	}
	if slices.IsSorted(sortStrs[:]) {
		t.Error("IsSorted: unsorted strs")
	}
	// IsSorted: true on sorted data.
	sorted := []int{1, 2, 3, 4, 5}
	if !slices.IsSorted(sorted) {
		t.Error("IsSorted: sorted ints")
	}
	sortedStrs := []string{"a", "b", "c"}
	if !slices.IsSorted(sortedStrs) {
		t.Error("IsSorted: sorted strs")
	}
}

func TestIsSortedFunc(t *testing.T) {
	// IsSortedFunc: false on unsorted data.
	compare := cmp.FuncFor[int]()
	if slices.IsSortedFunc(sortInts[:], compare) {
		t.Error("IsSortedFunc: unsorted ints")
	}
	// IsSortedFunc: true on sorted data.
	sorted := []int{1, 2, 3, 4, 5}
	if !slices.IsSortedFunc(sorted, compare) {
		t.Error("IsSortedFunc: sorted ints")
	}
}

func TestSort_Ints(t *testing.T) {
	s := slices.Clone(mem.System, sortInts[:])
	defer slices.Free(mem.System, s)
	slices.Sort(s)
	if !slices.IsSorted(s) {
		t.Error("Sort ints: not sorted")
	}
	if s[0] != -5467984 || s[12] != 9845 {
		t.Error("Sort ints: wrong values")
	}
}

func TestSort_Float64s(t *testing.T) {
	s := slices.Clone(mem.System, sortFloat64s[:])
	defer slices.Free(mem.System, s)
	slices.Sort(s)
	if !slices.IsSorted(s) {
		t.Error("Sort float64s: not sorted")
	}
	if s[0] != -959.7485 || s[14] != 9845.768 {
		t.Error("Sort float64s: wrong values")
	}
}

func TestSort_Strings(t *testing.T) {
	s := slices.Clone(mem.System, sortStrs[:])
	defer slices.Free(mem.System, s)
	slices.Sort(s)
	if !slices.IsSorted(s) {
		t.Error("Sort strings: not sorted")
	}
	if s[0] != "" || s[7] != "foo" {
		t.Error("Sort strings: wrong values")
	}
}

func TestSortFunc(t *testing.T) {
	// SortFunc (reverse order).
	s := slices.Clone(mem.System, sortInts[:])
	defer slices.Free(mem.System, s)
	slices.SortFunc(s, descInt)
	if !slices.IsSortedFunc(s, descInt) {
		t.Error("SortFunc ints: not sorted")
	}
	if s[0] != 9845 || s[12] != -5467984 {
		t.Error("SortFunc ints: wrong values")
	}
}

func TestSortFunc_Nil(t *testing.T) {
	// SortFunc with nil compare.
	type point struct{ x, y int }
	s := []point{{1, 2}, {3, 4}, {2, 3}}
	slices.SortFunc(s, nil)
	if !slices.IsSortedFunc(s, nil) {
		t.Error("SortFunc with nil: not sorted")
	}
	if s[0].x != 1 || s[0].y != 2 {
		t.Error("SortFunc with nil: wrong s[0]")
	}
	if s[1].x != 2 || s[1].y != 3 {
		t.Error("SortFunc with nil: wrong s[1]")
	}
	if s[2].x != 3 || s[2].y != 4 {
		t.Error("SortFunc with nil: wrong s[2]")
	}
}

func TestSortStableFunc_Ints(t *testing.T) {
	s := slices.Clone(mem.System, sortInts[:])
	defer slices.Free(mem.System, s)
	compare := cmp.FuncFor[int]()
	slices.SortStableFunc(s, compare)
	if !slices.IsSorted(s) {
		t.Error("SortStable ints: not sorted")
	}
	if s[0] != -5467984 || s[12] != 9845 {
		t.Error("SortStable ints: wrong values")
	}
}

func TestSortStableFunc_Float64s(t *testing.T) {
	s := slices.Clone(mem.System, sortFloat64s[:])
	defer slices.Free(mem.System, s)
	compare := cmp.FuncFor[float64]()
	slices.SortStableFunc(s, compare)
	if !slices.IsSorted(s) {
		t.Error("SortStable float64s: not sorted")
	}
	if s[0] != -959.7485 || s[14] != 9845.768 {
		t.Error("SortStable float64s: wrong values")
	}
}

func TestSortStableFunc_Strings(t *testing.T) {
	s := slices.Clone(mem.System, sortStrs[:])
	defer slices.Free(mem.System, s)
	compare := cmp.FuncFor[string]()
	slices.SortStableFunc(s, compare)
	if !slices.IsSorted(s) {
		t.Error("SortStable strings: not sorted")
	}
	if s[0] != "" || s[7] != "foo" {
		t.Error("SortStable strings: wrong values")
	}
}
