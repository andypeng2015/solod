package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/slices"
	"solod.dev/so/testing"
)

func TestMake(t *testing.T) {
	s := slices.Make[int](mem.System, 3)
	defer slices.Free(mem.System, s)
	s[0] = 11
	s[1] = 22
	s[2] = 33
	if len(s) != 3 || cap(s) != 3 {
		t.Error("Make failed")
	}
	if s[0] != 11 || s[1] != 22 || s[2] != 33 {
		t.Error("Make failed")
	}
}

func TestAppend(t *testing.T) {
	// Append within capacity.
	s := slices.MakeCap[int](mem.System, 0, 8)
	s = slices.Append(mem.System, s, 10, 20, 30)
	if len(s) != 3 || s[0] != 10 || s[1] != 20 || s[2] != 30 {
		t.Error("Append failed")
	}
	slices.Free(mem.System, s)
}

func TestAppend_Grow(t *testing.T) {
	// Append that triggers growth.
	s := slices.MakeCap[int](mem.System, 0, 2)
	s = slices.Append(mem.System, s, 1, 2)
	s = slices.Append(mem.System, s, 3, 4, 5)
	if len(s) != 5 || s[0] != 1 || s[4] != 5 {
		t.Error("Append grow failed")
	}
	slices.Free(mem.System, s)
}

func TestAppend_Nil(t *testing.T) {
	// Append to nil slice.
	var s []int
	s = slices.Append(mem.System, s, 10, 20, 30)
	if len(s) != 3 || s[0] != 10 || s[1] != 20 || s[2] != 30 {
		t.Error("Append to nil failed")
	}
	slices.Free(mem.System, s)
}

func TestExtend(t *testing.T) {
	// Extend from another slice.
	s := slices.MakeCap[int](mem.System, 0, 8)
	other := []int{100, 200, 300}
	s = slices.Extend(mem.System, s, other)
	if len(s) != 3 || s[0] != 100 || s[2] != 300 {
		t.Error("Extend failed")
	}
	slices.Free(mem.System, s)
}

func TestExtend_Nil(t *testing.T) {
	// Extend a nil slice.
	var s []int
	other := []int{10, 20, 30}
	s = slices.Extend(mem.System, s, other)
	if len(s) != 3 || s[0] != 10 || s[1] != 20 || s[2] != 30 {
		t.Error("Extend to nil failed")
	}
	slices.Free(mem.System, s)
}

func TestClone(t *testing.T) {
	s1 := []int{11, 22, 33}
	s2 := slices.Clone(mem.System, s1)
	defer slices.Free(mem.System, s2)
	s2[0] = 99
	if s1[0] != 11 || s2[0] != 99 {
		t.Error("Clone failed")
	}
}

func TestEqual(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{1, 2, 4}
	s4 := []int{1, 2}
	s5 := []int{}
	var s6 []int = nil
	if !slices.Equal(s1, s2) {
		t.Error("want s1 == s2")
	}
	if slices.Equal(s1, s3) {
		t.Error("want s1 != s3")
	}
	if slices.Equal(s1, s4) {
		t.Error("want s1 != s4")
	}
	if !slices.Equal(s5, s6) {
		t.Error("want empty and nil slices equal")
	}
}

func TestEqual_Strings(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"a", "b", "c"}
	s3 := []string{"a", "b", "d"}
	if !slices.Equal(s1, s2) {
		t.Error("want s1 == s2")
	}
	if slices.Equal(s1, s3) {
		t.Error("want s1 != s3")
	}
}

func TestEqual_Structs(t *testing.T) {
	type point struct {
		x, y int
	}
	s1 := []point{{1, 2}, {3, 4}}
	s2 := []point{{1, 2}, {3, 4}}
	s3 := []point{{1, 2}, {3, 5}}
	if !slices.Equal(s1, s2) {
		t.Error("want s1 == s2")
	}
	if slices.Equal(s1, s3) {
		t.Error("want s1 != s3")
	}
}

func TestIndex(t *testing.T) {
	ints := []int{10, 20, 30, 20}
	if slices.Index(ints, 20) != 1 {
		t.Error("Index failed")
	}
	if slices.Index(ints, 40) != -1 {
		t.Error("Index failed")
	}
	strs := []string{"a", "b", "c", "b"}
	if slices.Index(strs, "b") != 1 {
		t.Error("Index failed")
	}
	if slices.Index(strs, "d") != -1 {
		t.Error("Index failed")
	}
}

func TestContains(t *testing.T) {
	ints := []int{10, 20, 30, 20}
	if !slices.Contains(ints, 20) {
		t.Error("Contains failed")
	}
	if slices.Contains(ints, 40) {
		t.Error("Contains failed")
	}
	strs := []string{"a", "b", "c", "b"}
	if !slices.Contains(strs, "b") {
		t.Error("Contains failed")
	}
	if slices.Contains(strs, "d") {
		t.Error("Contains failed")
	}
}

func TestMinMax_Ints(t *testing.T) {
	ints := []int{3, 1, 4, 1, 5, 9}
	if slices.Min(ints) != 1 {
		t.Error("Min ints: wrong value")
	}
	if slices.Max(ints) != 9 {
		t.Error("Max ints: wrong value")
	}
}

func TestMinMax_Strings(t *testing.T) {
	strs := []string{"banana", "apple", "cherry"}
	if slices.Min(strs) != "apple" {
		t.Error("Min strings: wrong value")
	}
	if slices.Max(strs) != "cherry" {
		t.Error("Max strings: wrong value")
	}
}
