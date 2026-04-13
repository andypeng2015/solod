// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmp_test

import (
	"testing"
	"unsafe"

	"solod.dev/so/cmp"
	"solod.dev/so/fmt"
	"solod.dev/so/math"
)

var negzero = math.Copysign(0, -1)
var nonnilptr uintptr = uintptr(unsafe.Pointer(&negzero))
var nilptr uintptr = uintptr(unsafe.Pointer(nil))

var tests = []struct {
	x, y    any
	compare int
}{
	{1, 2, -1},
	{1, 1, 0},
	{2, 1, +1},
	{"a", "aa", -1},
	{"a", "a", 0},
	{"aa", "a", +1},
	{1.0, 1.1, -1},
	{1.1, 1.1, 0},
	{1.1, 1.0, +1},
	{math.Inf(1), math.Inf(1), 0},
	{math.Inf(-1), math.Inf(-1), 0},
	{math.Inf(1), 1.0, +1},
	{1.0, math.Inf(1), -1},
	{math.NaN(), math.NaN(), 0},
	{0.0, 0.0, 0},
	{negzero, negzero, 0},
	{negzero, 1.0, -1},
	{nilptr, nonnilptr, -1},
	{nonnilptr, nilptr, 1},
	{nonnilptr, nonnilptr, 0},
}

func TestLess(t *testing.T) {
	for _, test := range tests {
		var b bool
		switch test.x.(type) {
		case int:
			b = cmp.Less(test.x.(int), test.y.(int))
		case string:
			b = cmp.Less(test.x.(string), test.y.(string))
		case float64:
			b = cmp.Less(test.x.(float64), test.y.(float64))
		case uintptr:
			b = cmp.Less(test.x.(uintptr), test.y.(uintptr))
		}
		if b != (test.compare < 0) {
			t.Errorf("Less(%v, %v) == %t, want %t", test.x, test.y, b, test.compare < 0)
		}
	}
}

func TestCompare(t *testing.T) {
	for _, test := range tests {
		var c int
		switch test.x.(type) {
		case int:
			c = cmp.Compare(test.x.(int), test.y.(int))
		case string:
			c = cmp.Compare(test.x.(string), test.y.(string))
		case float64:
			c = cmp.Compare(test.x.(float64), test.y.(float64))
		case uintptr:
			c = cmp.Compare(test.x.(uintptr), test.y.(uintptr))
		}
		if c != test.compare {
			t.Errorf("Compare(%v, %v) == %d, want %d", test.x, test.y, c, test.compare)
		}
	}
}

func TestEqual(t *testing.T) {
	for _, test := range tests {
		var b bool
		switch test.x.(type) {
		case int:
			b = cmp.Equal(test.x.(int), test.y.(int))
		case string:
			b = cmp.Equal(test.x.(string), test.y.(string))
		case float64:
			b = cmp.Equal(test.x.(float64), test.y.(float64))
		case uintptr:
			b = cmp.Equal(test.x.(uintptr), test.y.(uintptr))
		}
		if b != (test.compare == 0) {
			t.Errorf("Equal(%v, %v) == %t, want %t", test.x, test.y, b, test.compare == 0)
		}
	}
}

func ExampleLess() {
	fmt.Printf("%t\n", cmp.Less(1, 2))
	fmt.Printf("%t\n", cmp.Less("a", "aa"))
	// Output:
	// true
	// true
}

func ExampleCompare() {
	fmt.Printf("%d\n", cmp.Compare(1, 2))
	fmt.Printf("%d\n", cmp.Compare("a", "aa"))
	fmt.Printf("%d\n", cmp.Compare(1.5, 1.5))
	// Output:
	// -1
	// -1
	// 0
}

func ExampleEqual() {
	fmt.Printf("%t\n", cmp.Equal(1, 1))
	fmt.Printf("%t\n", cmp.Equal("a", "aa"))
	fmt.Printf("%t\n", cmp.Equal(1.5, 1.5))
	// Output:
	// true
	// false
	// true
}
