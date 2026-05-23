// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netip

import "testing"

func TestUint128AddSub(t *testing.T) {
	const add1 = 1
	const sub1 = -1
	tests := []struct {
		in   uint128
		op   int // +1 or -1 to add vs subtract
		want uint128
	}{
		{uint128{0, 0}, add1, uint128{0, 1}},
		{uint128{0, 1}, add1, uint128{0, 2}},
		{uint128{1, 0}, add1, uint128{1, 1}},
		{uint128{0, ^uint64(0)}, add1, uint128{1, 0}},
		{uint128{^uint64(0), ^uint64(0)}, add1, uint128{0, 0}},

		{uint128{0, 0}, sub1, uint128{^uint64(0), ^uint64(0)}},
		{uint128{0, 1}, sub1, uint128{0, 0}},
		{uint128{0, 2}, sub1, uint128{0, 1}},
		{uint128{1, 0}, sub1, uint128{0, ^uint64(0)}},
		{uint128{1, 1}, sub1, uint128{1, 0}},
	}
	for _, tt := range tests {
		var got uint128
		switch tt.op {
		case add1:
			got = tt.in.addOne()
		case sub1:
			got = tt.in.subOne()
		default:
			panic("bogus op")
		}
		if got != tt.want {
			t.Errorf("%v add %d = %v; want %v", tt.in, tt.op, got, tt.want)
		}
	}
}

func TestMask6(t *testing.T) {
	tests := []struct {
		n    int
		want uint128
	}{
		{0, uint128{0, 0}},
		{1, uint128{0x8000000000000000, 0}},
		{32, uint128{0xFFFFFFFF00000000, 0}},
		{63, uint128{0xFFFFFFFFFFFFFFFE, 0}},
		{64, uint128{^uint64(0), 0}},
		{65, uint128{^uint64(0), 0x8000000000000000}},
		{96, uint128{^uint64(0), 0xFFFFFFFF00000000}},
		{127, uint128{^uint64(0), 0xFFFFFFFFFFFFFFFE}},
		{128, uint128{^uint64(0), ^uint64(0)}},
	}
	for _, tt := range tests {
		got := mask6(tt.n)
		if got != tt.want {
			t.Errorf("mask6(%d) = %v; want %v", tt.n, got, tt.want)
		}
	}
}
