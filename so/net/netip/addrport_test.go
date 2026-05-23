// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netip

import (
	"slices"
	"testing"
)

func TestAddrPortCompare(t *testing.T) {
	tests := []struct {
		a, b AddrPort
		want int
	}{
		{AddrPort{}, AddrPort{}, 0},
		{AddrPort{}, MustParseAddrPort("1.2.3.4:80"), -1},

		{MustParseAddrPort("1.2.3.4:80"), MustParseAddrPort("1.2.3.4:80"), 0},
		{MustParseAddrPort("[::1]:80"), MustParseAddrPort("[::1]:80"), 0},

		{MustParseAddrPort("1.2.3.4:80"), MustParseAddrPort("2.3.4.5:22"), -1},
		{MustParseAddrPort("[::1]:80"), MustParseAddrPort("[::2]:22"), -1},

		{MustParseAddrPort("1.2.3.4:80"), MustParseAddrPort("1.2.3.4:443"), -1},
		{MustParseAddrPort("[::1]:80"), MustParseAddrPort("[::1]:443"), -1},

		{MustParseAddrPort("1.2.3.4:80"), MustParseAddrPort("[0102:0304::0]:80"), -1},
	}
	var buf [64]byte
	for _, tt := range tests {
		a := tt.a.String(buf[:])
		b := tt.b.String(buf[:])
		got := tt.a.Compare(tt.b)
		if got != tt.want {
			t.Errorf("Compare(%q, %q) = %v; want %v", a, b, got, tt.want)
		}

		// Also check inverse.
		if got == tt.want {
			got2 := tt.b.Compare(tt.a)
			if want2 := -1 * tt.want; got2 != want2 {
				t.Errorf("Compare(%q, %q) was correctly %v, but Compare(%q, %q) was %v", b, a, got, b, a, got2)
			}
		}
	}

	// And just sort.
	values := []AddrPort{
		MustParseAddrPort("[::1]:80"),
		MustParseAddrPort("[::2]:80"),
		AddrPort{},
		MustParseAddrPort("1.2.3.4:443"),
		MustParseAddrPort("8.8.8.8:8080"),
		MustParseAddrPort("[::1%foo]:1024"),
	}
	sorted := slices.Clone(values)
	slices.SortFunc(sorted, AddrPort.Compare)
	want := []int{2, 3, 4, 0, 5, 1} // indices of values in sorted order
	for i, v := range sorted {
		if v != values[want[i]] {
			gots := v.String(buf[:])
			wants := values[want[i]].String(buf[:])
			t.Errorf("unexpected sort at index %d: got %q, want %q", i, gots, wants)
		}
	}
}

func TestAddrPortString(t *testing.T) {
	tests := []struct {
		ipp  AddrPort
		want string
	}{
		{MustParseAddrPort("127.0.0.1:80"), "127.0.0.1:80"},
		{MustParseAddrPort("[0000::0]:8080"), "[::]:8080"},
		{MustParseAddrPort("[FFFF::1]:8080"), "[ffff::1]:8080"},
		{AddrPort{}, "invalid AddrPort"},
		{AddrPortFrom(Addr{}, 80), "invalid AddrPort"},
	}

	for _, tt := range tests {
		var buf [64]byte
		if got := tt.ipp.String(buf[:]); got != tt.want {
			t.Errorf("(%#v).String() = %q want %q", tt.ipp, got, tt.want)
		}
	}
}
