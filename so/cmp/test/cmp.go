package main

import (
	"solod.dev/so/cmp"
	"solod.dev/so/testing"
)

func TestCompare(t *testing.T) {
	a, b := 11, 22
	if cmp.Compare(a, b) >= 0 {
		t.Error("Compare(11, 22) >= 0")
	}
	if cmp.Compare(a, a) != 0 {
		t.Error("Compare(11, 11) != 0")
	}

	s1, s2 := "hello", "world"
	if cmp.Compare(s1, s2) >= 0 {
		t.Error("Compare(hello, world) >= 0")
	}
	if cmp.Compare(s1, s1) != 0 {
		t.Error("Compare(hello, hello) != 0")
	}
}

func TestEqual(t *testing.T) {
	a, b := 11, 22
	if cmp.Equal(a, b) {
		t.Error("Equal(11, 22) = true")
	}
	if !cmp.Equal(a, a) {
		t.Error("Equal(11, 11) = false")
	}

	s1, s2 := "hello", "world"
	if cmp.Equal(s1, s2) {
		t.Error("Equal(hello, world) = true")
	}
	if !cmp.Equal(s1, s1) {
		t.Error("Equal(hello, hello) = false")
	}
}

func TestLess(t *testing.T) {
	a, b := 11, 22
	if !cmp.Less(a, b) {
		t.Error("Less(11, 22) = false")
	}
	if cmp.Less(b, a) {
		t.Error("Less(22, 11) = true")
	}

	s1, s2 := "hello", "world"
	if !cmp.Less(s1, s2) {
		t.Error("Less(hello, world) = false")
	}
	if cmp.Less(s2, s1) {
		t.Error("Less(world, hello) = true")
	}
}
