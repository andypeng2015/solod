// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netip

import (
	"net"
	"reflect"
	"slices"
	"testing"

	"solod.dev/so/bytes"
)

func TestParseAddr(t *testing.T) {
	var validIPs = []struct {
		in      string
		ip      Addr   // output of ParseAddr()
		str     string // output of String(). If "", use in.
		wantErr error
	}{
		// Basic zero IPv4 address.
		{
			in: "0.0.0.0",
			ip: makeAddr(make128(0, 0xffff00000000), z4),
		},
		// Basic non-zero IPv4 address.
		{
			in: "192.168.140.255",
			ip: makeAddr(make128(0, 0xffffc0a88cff), z4),
		},
		// IPv4 address in windows-style "print all the digits" form.
		{
			in:      "010.000.015.001",
			wantErr: ErrIPv4,
		},
		// IPv4 address with a silly amount of leading zeros.
		{
			in:      "000001.00000002.00000003.000000004",
			wantErr: ErrIPv4,
		},
		// 4-in-6 with octet with leading zero
		{
			in:      "::ffff:1.2.03.4",
			wantErr: ErrIPv4,
		},
		// 4-in-6 with octet with unexpected character
		{
			in:      "::ffff:1.2.3.z",
			wantErr: ErrIPv4,
		},
		// Basic zero IPv6 address.
		{
			in: "::",
			ip: makeAddr(make128(0, 0), z6),
		},
		// Localhost IPv6.
		{
			in: "::1",
			ip: makeAddr(make128(0, 1), z6),
		},
		// Fully expanded IPv6 address.
		{
			in: "fd7a:115c:a1e0:ab12:4843:cd96:626b:430b",
			ip: makeAddr(make128(0xfd7a115ca1e0ab12, 0x4843cd96626b430b), z6),
		},
		// IPv6 with elided fields in the middle.
		{
			in: "fd7a:115c::626b:430b",
			ip: makeAddr(make128(0xfd7a115c00000000, 0x00000000626b430b), z6),
		},
		// IPv6 with elided fields at the end.
		{
			in: "fd7a:115c:a1e0:ab12:4843:cd96::",
			ip: makeAddr(make128(0xfd7a115ca1e0ab12, 0x4843cd9600000000), z6),
		},
		// IPv6 with single elided field at the end.
		{
			in:  "fd7a:115c:a1e0:ab12:4843:cd96:626b::",
			ip:  makeAddr(make128(0xfd7a115ca1e0ab12, 0x4843cd96626b0000), z6),
			str: "fd7a:115c:a1e0:ab12:4843:cd96:626b:0",
		},
		// IPv6 with single elided field in the middle.
		{
			in:  "fd7a:115c:a1e0::4843:cd96:626b:430b",
			ip:  makeAddr(make128(0xfd7a115ca1e00000, 0x4843cd96626b430b), z6),
			str: "fd7a:115c:a1e0:0:4843:cd96:626b:430b",
		},
		// IPv6 with the trailing 32 bits written as IPv4 dotted decimal. (4in6)
		{
			in:  "::ffff:192.168.140.255",
			ip:  makeAddr(make128(0, 0x0000ffffc0a88cff), z6),
			str: "::ffff:192.168.140.255",
		},
		// IPv6 with a zone specifier.
		{
			in:  "fd7a:115c:a1e0:ab12:4843:cd96:626b:430b%eth0",
			ip:  makeAddrZone(make128(0xfd7a115ca1e0ab12, 0x4843cd96626b430b), z6, 10),
			str: "fd7a:115c:a1e0:ab12:4843:cd96:626b:430b%10",
		},
		// IPv6 with dotted decimal and zone specifier.
		{
			in:  "1:2::ffff:192.168.140.255%eth1",
			ip:  makeAddrZone(make128(0x0001000200000000, 0x0000ffffc0a88cff), z6, 11),
			str: "1:2::ffff:c0a8:8cff%11",
		},
		// 4-in-6 with zone
		{
			in:  "::ffff:192.168.140.255%eth1",
			ip:  makeAddrZone(make128(0, 0x0000ffffc0a88cff), z6, 11),
			str: "::ffff:192.168.140.255%11",
		},
		// IPv6 with capital letters.
		{
			in:  "FD9E:1A04:F01D::1",
			ip:  makeAddr(make128(0xfd9e1a04f01d0000, 0x1), z6),
			str: "fd9e:1a04:f01d::1",
		},
	}

	for _, test := range validIPs {
		t.Run(test.in, func(t *testing.T) {
			var buf [64]byte

			got, err := ParseAddr(test.in)
			if err != nil {
				if err == test.wantErr {
					return
				}
				t.Fatalf("wanted error %q; got %q", test.wantErr, err)
			}
			if test.wantErr != nil {
				t.Fatalf("wanted error %q; got none", test.wantErr)
			}
			if got != test.ip {
				t.Errorf("got %#v, want %#v", got, test.ip)
			}

			// Check that ParseAddr is a pure function.
			got2, err := ParseAddr(test.in)
			if err != nil {
				t.Fatal(err)
			}
			if got != got2 {
				t.Errorf("ParseAddr(%q) got 2 different results: %#v, %#v", test.in, got, got2)
			}

			// Check that ParseAddr(ip.String()) is the identity function.

			s := got.String(buf[:])
			got3, err := ParseAddr(s)
			if err != nil {
				t.Fatal(err)
			}
			if got != got3 {
				t.Errorf("ParseAddr(%q) != ParseAddr(ParseIP(%q).String()). Got %#v, want %#v", test.in, test.in, got3, got)
			}

			// Check that the parsed IP formats as expected.
			s = got.String(buf[:])
			wants := test.str
			if wants == "" {
				wants = test.in
			}
			if s != wants {
				t.Errorf("ParseAddr(%q).String() got %q, want %q", test.in, s, wants)
			}
		})
	}

	var invalidIPs = []string{
		// Empty string
		"",
		// Garbage non-IP
		"bad",
		// Single number. Some parsers accept this as an IPv4 address in
		// big-endian uint32 form, but we don't.
		"1234",
		// IPv4 with a zone specifier
		"1.2.3.4%eth0",
		// IPv4 field must have at least one digit
		".1.2.3",
		"1.2.3.",
		"1..2.3",
		// IPv4 address too long
		"1.2.3.4.5",
		// IPv4 in dotted octal form
		"0300.0250.0214.0377",
		// IPv4 in dotted hex form
		"0xc0.0xa8.0x8c.0xff",
		// IPv4 in class B form
		"192.168.12345",
		// IPv4 in class B form, with a small enough number to be
		// parseable as a regular dotted decimal field.
		"127.0.1",
		// IPv4 in class A form
		"192.1234567",
		// IPv4 in class A form, with a small enough number to be
		// parseable as a regular dotted decimal field.
		"127.1",
		// IPv4 field has value >255
		"192.168.300.1",
		// IPv4 with too many fields
		"192.168.0.1.5.6",
		// IPv6 with not enough fields
		"1:2:3:4:5:6:7",
		// IPv6 with too many fields
		"1:2:3:4:5:6:7:8:9",
		// IPv6 with 8 fields and a :: expander
		"1:2:3:4::5:6:7:8",
		// IPv6 with a field bigger than 2b
		"fe801::1",
		// IPv6 with non-hex values in field
		"fe80:tail:scal:e::",
		// IPv6 with a zone delimiter but no zone.
		"fe80::1%",
		// IPv6 (without ellipsis) with too many fields for trailing embedded IPv4.
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:192.168.140.255",
		// IPv6 (with ellipsis) with too many fields for trailing embedded IPv4.
		"ffff::ffff:ffff:ffff:ffff:ffff:ffff:192.168.140.255",
		// IPv6 with invalid embedded IPv4.
		"::ffff:192.168.140.bad",
		// IPv6 with multiple ellipsis ::.
		"fe80::1::1",
		// IPv6 with invalid non hex/colon character.
		"fe80:1?:1",
		// IPv6 with truncated bytes after single colon.
		"fe80:",
		// IPv6 with 5 zeros in last group
		"0:0:0:0:0:ffff:0:00000",
		// IPv6 with 5 zeros in one group and embedded IPv4
		"0:0:0:0:00000:ffff:127.1.2.3",
	}

	for _, s := range invalidIPs {
		t.Run(s, func(t *testing.T) {
			got, err := ParseAddr(s)
			if err == nil {
				t.Errorf("ParseAddr(%q) = %#v, want error", s, got)
			}
		})
	}
}

func TestAddrFromSlice(t *testing.T) {
	tests := []struct {
		ip       []byte
		wantAddr Addr
		wantOK   bool
	}{
		{
			ip:       []byte{10, 0, 0, 1},
			wantAddr: MustParseAddr("10.0.0.1"),
			wantOK:   true,
		},
		{
			ip:       []byte{0xfe, 0x80, 15: 0x01},
			wantAddr: MustParseAddr("fe80::01"),
			wantOK:   true,
		},
		{
			ip:       []byte{0, 1, 2},
			wantAddr: Addr{},
			wantOK:   false,
		},
		{
			ip:       nil,
			wantAddr: Addr{},
			wantOK:   false,
		},
	}
	for _, tt := range tests {
		addr := AddrFromSlice(tt.ip)
		if addr != tt.wantAddr {
			t.Errorf("AddrFromSlice(%#v) = %#v, want %#v", tt.ip, addr, tt.wantAddr)
		}
	}
}

func TestIPv4Constructors(t *testing.T) {
	if AddrFrom4([4]byte{1, 2, 3, 4}) != MustParseAddr("1.2.3.4") {
		t.Errorf("don't match")
	}
}

func TestAddrAppendText(t *testing.T) {
	tests := []struct {
		ip   Addr
		want string
	}{
		{Addr{}, ""}, // zero IP
		{MustParseAddr("1.2.3.4"), "1.2.3.4"},
		{MustParseAddr("fd7a:115c:a1e0:ab12:4843:cd96:626b:430b"), "fd7a:115c:a1e0:ab12:4843:cd96:626b:430b"},
		{MustParseAddr("::ffff:192.168.140.255"), "::ffff:192.168.140.255"},
		{MustParseAddr("::ffff:192.168.140.255%en0"), "::ffff:192.168.140.255%10"},
	}
	for i, tc := range tests {
		ip := tc.ip

		mtAppend := make([]byte, 4, 32)
		mtAppend, err := ip.AppendText(mtAppend)
		mtAppend = mtAppend[4:]
		if err != nil {
			t.Fatal(err)
		}
		if string(mtAppend) != tc.want {
			t.Errorf("%d. for (%v) AppendText = %q; want %q", i, ip, mtAppend, tc.want)
		}
	}
}

func TestAddrFrom16(t *testing.T) {
	tests := []struct {
		name string
		in   [16]byte
		want Addr
	}{
		{
			name: "v6-raw",
			in:   [...]byte{15: 1},
			want: makeAddr(make128(0, 1), z6),
		},
		{
			name: "v4-raw",
			in:   [...]byte{10: 0xff, 11: 0xff, 12: 1, 13: 2, 14: 3, 15: 4},
			want: makeAddr(make128(0, 0xffff01020304), z6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddrFrom16(tt.in)
			if got != tt.want {
				t.Errorf("got %#v; want %#v", got, tt.want)
			}
		})
	}
}

func TestIPProperties(t *testing.T) {
	var (
		nilIP Addr

		unicast4           = MustParseAddr("192.0.2.1")
		unicast6           = MustParseAddr("2001:db8::1")
		unicastZone6       = MustParseAddr("2001:db8::1%eth0")
		unicast6Unassigned = MustParseAddr("4000::1") // not in 2000::/3.

		multicast4     = MustParseAddr("224.0.0.1")
		multicast6     = MustParseAddr("ff02::1")
		multicastZone6 = MustParseAddr("ff02::1%eth0")

		llu4     = MustParseAddr("169.254.0.1")
		llu6     = MustParseAddr("fe80::1")
		llu6Last = MustParseAddr("febf:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
		lluZone6 = MustParseAddr("fe80::1%eth0")

		loopback4 = MustParseAddr("127.0.0.1")

		ilm6     = MustParseAddr("ff01::1")
		ilmZone6 = MustParseAddr("ff01::1%eth0")

		private4a        = MustParseAddr("10.0.0.1")
		private4b        = MustParseAddr("172.16.0.1")
		private4c        = MustParseAddr("192.168.1.1")
		private6         = MustParseAddr("fd00::1")
		private6mapped4a = MustParseAddr("::ffff:10.0.0.1")
		private6mapped4b = MustParseAddr("::ffff:172.16.0.1")
		private6mapped4c = MustParseAddr("::ffff:192.168.1.1")
	)

	var a16 [16]byte
	tests := []struct {
		name                    string
		ip                      Addr
		globalUnicast           bool
		interfaceLocalMulticast bool
		linkLocalMulticast      bool
		linkLocalUnicast        bool
		loopback                bool
		multicast               bool
		private                 bool
		unspecified             bool
	}{
		{
			name: "nil",
			ip:   nilIP,
		},
		{
			name:          "unicast v4Addr",
			ip:            unicast4,
			globalUnicast: true,
		},
		{
			name:          "unicast v6 mapped v4Addr",
			ip:            AddrFrom16(unicast4.As16(a16)),
			globalUnicast: true,
		},
		{
			name:          "unicast v6Addr",
			ip:            unicast6,
			globalUnicast: true,
		},
		{
			name:          "unicast v6AddrZone",
			ip:            unicastZone6,
			globalUnicast: true,
		},
		{
			name:          "unicast v6Addr unassigned",
			ip:            unicast6Unassigned,
			globalUnicast: true,
		},
		{
			name:               "multicast v4Addr",
			ip:                 multicast4,
			linkLocalMulticast: true,
			multicast:          true,
		},
		{
			name:               "multicast v6 mapped v4Addr",
			ip:                 AddrFrom16(multicast4.As16(a16)),
			linkLocalMulticast: true,
			multicast:          true,
		},
		{
			name:               "multicast v6Addr",
			ip:                 multicast6,
			linkLocalMulticast: true,
			multicast:          true,
		},
		{
			name:               "multicast v6AddrZone",
			ip:                 multicastZone6,
			linkLocalMulticast: true,
			multicast:          true,
		},
		{
			name:             "link-local unicast v4Addr",
			ip:               llu4,
			linkLocalUnicast: true,
		},
		{
			name:             "link-local unicast v6 mapped v4Addr",
			ip:               AddrFrom16(llu4.As16(a16)),
			linkLocalUnicast: true,
		},
		{
			name:             "link-local unicast v6Addr",
			ip:               llu6,
			linkLocalUnicast: true,
		},
		{
			name:             "link-local unicast v6Addr upper bound",
			ip:               llu6Last,
			linkLocalUnicast: true,
		},
		{
			name:             "link-local unicast v6AddrZone",
			ip:               lluZone6,
			linkLocalUnicast: true,
		},
		{
			name:     "loopback v4Addr",
			ip:       loopback4,
			loopback: true,
		},
		{
			name:     "loopback v6Addr",
			ip:       IPv6Loopback(),
			loopback: true,
		},
		{
			name:     "loopback v6 mapped v4Addr",
			ip:       AddrFrom16(IPv6Loopback().As16(a16)),
			loopback: true,
		},
		{
			name:                    "interface-local multicast v6Addr",
			ip:                      ilm6,
			interfaceLocalMulticast: true,
			multicast:               true,
		},
		{
			name:                    "interface-local multicast v6AddrZone",
			ip:                      ilmZone6,
			interfaceLocalMulticast: true,
			multicast:               true,
		},
		{
			name:          "private v4Addr 10/8",
			ip:            private4a,
			globalUnicast: true,
			private:       true,
		},
		{
			name:          "private v4Addr 172.16/12",
			ip:            private4b,
			globalUnicast: true,
			private:       true,
		},
		{
			name:          "private v4Addr 192.168/16",
			ip:            private4c,
			globalUnicast: true,
			private:       true,
		},
		{
			name:          "private v6Addr",
			ip:            private6,
			globalUnicast: true,
			private:       true,
		},
		{
			name:          "private v6 mapped v4Addr 10/8",
			ip:            private6mapped4a,
			globalUnicast: true,
			private:       true,
		},
		{
			name:          "private v6 mapped v4Addr 172.16/12",
			ip:            private6mapped4b,
			globalUnicast: true,
			private:       true,
		},
		{
			name:          "private v6 mapped v4Addr 192.168/16",
			ip:            private6mapped4c,
			globalUnicast: true,
			private:       true,
		},
		{
			name:        "unspecified v4Addr",
			ip:          IPv4Unspecified(),
			unspecified: true,
		},
		{
			name:        "unspecified v6Addr",
			ip:          IPv6Unspecified(),
			unspecified: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gu := tt.ip.IsGlobalUnicast()
			if gu != tt.globalUnicast {
				t.Errorf("IsGlobalUnicast(%v) = %v; want %v", tt.ip, gu, tt.globalUnicast)
			}

			ilm := tt.ip.IsInterfaceLocalMulticast()
			if ilm != tt.interfaceLocalMulticast {
				t.Errorf("IsInterfaceLocalMulticast(%v) = %v; want %v", tt.ip, ilm, tt.interfaceLocalMulticast)
			}

			llu := tt.ip.IsLinkLocalUnicast()
			if llu != tt.linkLocalUnicast {
				t.Errorf("IsLinkLocalUnicast(%v) = %v; want %v", tt.ip, llu, tt.linkLocalUnicast)
			}

			llm := tt.ip.IsLinkLocalMulticast()
			if llm != tt.linkLocalMulticast {
				t.Errorf("IsLinkLocalMulticast(%v) = %v; want %v", tt.ip, llm, tt.linkLocalMulticast)
			}

			lo := tt.ip.IsLoopback()
			if lo != tt.loopback {
				t.Errorf("IsLoopback(%v) = %v; want %v", tt.ip, lo, tt.loopback)
			}

			multicast := tt.ip.IsMulticast()
			if multicast != tt.multicast {
				t.Errorf("IsMulticast(%v) = %v; want %v", tt.ip, multicast, tt.multicast)
			}

			private := tt.ip.IsPrivate()
			if private != tt.private {
				t.Errorf("IsPrivate(%v) = %v; want %v", tt.ip, private, tt.private)
			}

			unspecified := tt.ip.IsUnspecified()
			if unspecified != tt.unspecified {
				t.Errorf("IsUnspecified(%v) = %v; want %v", tt.ip, unspecified, tt.unspecified)
			}
		})
	}
}

func TestAddrWellKnown(t *testing.T) {
	tests := []struct {
		name string
		ip   Addr
		std  net.IP
	}{
		{
			name: "IPv4 unspecified",
			ip:   IPv4Unspecified(),
			std:  net.IPv4zero,
		},
		{
			name: "IPv6 link-local all nodes",
			ip:   IPv6LinkLocalAllNodes(),
			std:  net.IPv6linklocalallnodes,
		},
		{
			name: "IPv6 link-local all routers",
			ip:   IPv6LinkLocalAllRouters(),
			std:  net.IPv6linklocalallrouters,
		},
		{
			name: "IPv6 loopback",
			ip:   IPv6Loopback(),
			std:  net.IPv6loopback,
		},
		{
			name: "IPv6 unspecified",
			ip:   IPv6Unspecified(),
			std:  net.IPv6unspecified,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf [64]byte
			want := tt.std.String()
			got := tt.ip.String(buf[:])

			if got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		})
	}
}

func TestAddrLessCompare(t *testing.T) {
	tests := []struct {
		a, b Addr
		want bool
	}{
		{Addr{}, Addr{}, false},
		{Addr{}, MustParseAddr("1.2.3.4"), true},
		{MustParseAddr("1.2.3.4"), Addr{}, false},

		{MustParseAddr("1.2.3.4"), MustParseAddr("0102:0304::0"), true},
		{MustParseAddr("0102:0304::0"), MustParseAddr("1.2.3.4"), false},
		{MustParseAddr("1.2.3.4"), MustParseAddr("1.2.3.4"), false},

		{MustParseAddr("::1"), MustParseAddr("::2"), true},
		{MustParseAddr("::1"), MustParseAddr("::1%foo"), true},
		{MustParseAddr("::1%foo"), MustParseAddr("::2"), true},
		{MustParseAddr("::2"), MustParseAddr("::3"), true},

		{MustParseAddr("::"), MustParseAddr("0.0.0.0"), false},
		{MustParseAddr("0.0.0.0"), MustParseAddr("::"), true},

		{MustParseAddr("::1%a"), MustParseAddr("::1%b"), true},
		{MustParseAddr("::1%a"), MustParseAddr("::1%a"), false},
		{MustParseAddr("::1%b"), MustParseAddr("::1%a"), false},

		// For Issue 68113, verify that an IPv4 address and a
		// v4-mapped-IPv6 address differing only in their zone
		// pointer are unequal via all three of
		// ==/Compare/reflect.DeepEqual. In Go 1.22 and
		// earlier, these were accidentally equal via
		// DeepEqual due to their zone pointers (z) differing
		// but pointing to identical structures.
		{MustParseAddr("::ffff:11.1.1.12"), MustParseAddr("11.1.1.12"), false},
	}
	for _, tt := range tests {
		var buf [64]byte
		var a = tt.a.String(buf[:])
		var b = tt.b.String(buf[:])
		got := tt.a.Less(tt.b)
		if got != tt.want {
			t.Errorf("Less(%q, %q) = %v; want %v", a, b, got, tt.want)
		}
		cmp := tt.a.Compare(tt.b)
		if got && cmp != -1 {
			t.Errorf("Less(%q, %q) = true, but Compare = %v (not -1)", a, b, cmp)
		}
		if cmp < -1 || cmp > 1 {
			t.Errorf("bogus Compare return value %v", cmp)
		}
		if cmp == 0 && tt.a != tt.b {
			t.Errorf("Compare(%q, %q) = 0; but not equal", a, b)
		}
		if cmp == 1 && !tt.b.Less(tt.a) {
			t.Errorf("Compare(%q, %q) = 1; but b.Less(a) isn't true", a, b)
		}

		// Also check inverse.
		if got == tt.want && got {
			var buf2 [64]byte
			a2 := tt.a.String(buf2[:])
			b2 := tt.b.String(buf2[:])
			got2 := tt.b.Less(tt.a)
			if got2 {
				t.Errorf("Less(%q, %q) was correctly %v, but so was Less(%q, %q)", a, b, got, b2, a2)
			}
		}

		// Also check reflect.DeepEqual. See issue 68113.
		deepEq := reflect.DeepEqual(tt.a, tt.b)
		if (cmp == 0) != deepEq {
			t.Errorf("%q and %q differ in == (%v) vs reflect.DeepEqual (%v)", a, b, cmp == 0, deepEq)
		}
	}

	// And just sort.
	values := []Addr{
		MustParseAddr("::1"),
		MustParseAddr("::2"),
		Addr{},
		MustParseAddr("1.2.3.4"),
		MustParseAddr("8.8.8.8"),
		MustParseAddr("::1%foo"),
	}
	sorted := slices.Clone(values)
	slices.SortFunc(sorted, Addr.Compare)
	want := []int{2, 3, 4, 0, 5, 1} // indices of values in sorted order
	for i, v := range sorted {
		if v != values[want[i]] {
			t.Errorf("unexpected sort at index %d: got %v, want %v", i, v, values[want[i]])
		}
	}
}

func TestIs4AndIs6(t *testing.T) {
	tests := []struct {
		ip  Addr
		is4 bool
		is6 bool
	}{
		{Addr{}, false, false},
		{MustParseAddr("1.2.3.4"), true, false},
		{MustParseAddr("127.0.0.2"), true, false},
		{MustParseAddr("::1"), false, true},
		{MustParseAddr("::ffff:192.0.2.128"), false, true},
		{MustParseAddr("::fffe:c000:0280"), false, true},
		{MustParseAddr("::1%eth0"), false, true},
	}
	for _, tt := range tests {
		var buf [64]byte
		ip := tt.ip.String(buf[:])

		got4 := tt.ip.Is4()
		if got4 != tt.is4 {
			t.Errorf("Is4(%q) = %v; want %v", ip, got4, tt.is4)
		}

		got6 := tt.ip.Is6()
		if got6 != tt.is6 {
			t.Errorf("Is6(%q) = %v; want %v", ip, got6, tt.is6)
		}
	}
}

func TestIs4In6(t *testing.T) {
	tests := []struct {
		ip        Addr
		want      bool
		wantUnmap Addr
	}{
		{Addr{}, false, Addr{}},
		{MustParseAddr("::ffff:c000:0280"), true, MustParseAddr("192.0.2.128")},
		{MustParseAddr("::ffff:192.0.2.128"), true, MustParseAddr("192.0.2.128")},
		{MustParseAddr("::ffff:192.0.2.128%eth0"), true, MustParseAddr("192.0.2.128")},
		{MustParseAddr("::fffe:c000:0280"), false, MustParseAddr("::fffe:c000:0280")},
		{MustParseAddr("::ffff:127.1.2.3"), true, MustParseAddr("127.1.2.3")},
		{MustParseAddr("::ffff:7f01:0203"), true, MustParseAddr("127.1.2.3")},
		{MustParseAddr("0:0:0:0:0000:ffff:127.1.2.3"), true, MustParseAddr("127.1.2.3")},
		{MustParseAddr("0:0:0:0::ffff:127.1.2.3"), true, MustParseAddr("127.1.2.3")},
		{MustParseAddr("::1"), false, MustParseAddr("::1")},
		{MustParseAddr("1.2.3.4"), false, MustParseAddr("1.2.3.4")},
	}
	for _, tt := range tests {
		var buf [64]byte
		ip := tt.ip.String(buf[:])

		got := tt.ip.Is4In6()
		if got != tt.want {
			t.Errorf("Is4In6(%q) = %v; want %v", ip, got, tt.want)
		}
		u := tt.ip.Unmap()
		if u != tt.wantUnmap {
			t.Errorf("Unmap(%q) = %v; want %v", ip, u, tt.wantUnmap)
		}
	}
}

func TestAs4(t *testing.T) {
	var a16 [16]byte
	tests := []struct {
		ip        Addr
		want      [4]byte
		wantPanic bool
	}{
		{
			ip:   MustParseAddr("1.2.3.4"),
			want: [4]byte{1, 2, 3, 4},
		},
		{
			ip:   AddrFrom16(MustParseAddr("1.2.3.4").As16(a16)), // IPv4-in-IPv6
			want: [4]byte{1, 2, 3, 4},
		},
		{
			ip:   MustParseAddr("0.0.0.0"),
			want: [4]byte{0, 0, 0, 0},
		},
		{
			ip:        Addr{},
			wantPanic: true,
		},
		{
			ip:        MustParseAddr("::1"),
			wantPanic: true,
		},
	}
	as4 := func(ip Addr) (v [4]byte, gotPanic bool) {
		defer func() {
			if recover() != nil {
				gotPanic = true
				return
			}
		}()
		var a4 [4]byte
		v = ip.As4(a4)
		return
	}
	for i, tt := range tests {
		var buf [64]byte
		ip := tt.ip.String(buf[:])

		got, gotPanic := as4(tt.ip)
		if gotPanic != tt.wantPanic {
			t.Errorf("%d. panic on %v = %v; want %v", i, ip, gotPanic, tt.wantPanic)
			continue
		}
		if got != tt.want {
			t.Errorf("%d. %v = %v; want %v", i, ip, got, tt.want)
		}
	}
}

func TestAsSlice(t *testing.T) {
	tests := []struct {
		in   Addr
		want []byte
	}{
		{in: Addr{}, want: nil},
		{in: MustParseAddr("1.2.3.4"), want: []byte{1, 2, 3, 4}},
		{in: MustParseAddr("ffff::1"), want: []byte{0xff, 0xff, 15: 1}},
	}

	for _, test := range tests {
		var buf [64]byte
		got := test.in.AsSlice(buf[:])
		if !bytes.Equal(got, test.want) {
			t.Errorf("%v.AsSlice() = %v want %v", test.in, got, test.want)
		}
	}
}

func TestPrefixValid(t *testing.T) {
	v4 := MustParseAddr("1.2.3.4")
	v6 := MustParseAddr("::1")
	tests := []struct {
		ipp  Prefix
		want bool
	}{
		{PrefixFrom(v4, -2), false},
		{PrefixFrom(v4, -1), false},
		{PrefixFrom(v4, 0), true},
		{PrefixFrom(v4, 32), true},
		{PrefixFrom(v4, 33), false},

		{PrefixFrom(v6, -2), false},
		{PrefixFrom(v6, -1), false},
		{PrefixFrom(v6, 0), true},
		{PrefixFrom(v6, 32), true},
		{PrefixFrom(v6, 128), true},
		{PrefixFrom(v6, 129), false},

		{PrefixFrom(Addr{}, -2), false},
		{PrefixFrom(Addr{}, -1), false},
		{PrefixFrom(Addr{}, 0), false},
		{PrefixFrom(Addr{}, 32), false},
		{PrefixFrom(Addr{}, 128), false},
	}
	for _, tt := range tests {
		got := tt.ipp.IsValid()
		if got != tt.want {
			t.Errorf("(%v).IsValid() = %v want %v", tt.ipp, got, tt.want)
		}

		// Test that there is only one invalid Prefix representation per Addr.
		invalid := PrefixFrom(tt.ipp.Addr(), -1)
		if !got && tt.ipp != invalid {
			t.Errorf("(%v == %v) = false, want true", tt.ipp, invalid)
		}
	}
}

var nextPrevTests = []struct {
	ip   Addr
	next Addr
	prev Addr
}{
	{MustParseAddr("10.0.0.1"), MustParseAddr("10.0.0.2"), MustParseAddr("10.0.0.0")},
	{MustParseAddr("10.0.0.255"), MustParseAddr("10.0.1.0"), MustParseAddr("10.0.0.254")},
	{MustParseAddr("127.0.0.1"), MustParseAddr("127.0.0.2"), MustParseAddr("127.0.0.0")},
	{MustParseAddr("254.255.255.255"), MustParseAddr("255.0.0.0"), MustParseAddr("254.255.255.254")},
	{MustParseAddr("255.255.255.255"), Addr{}, MustParseAddr("255.255.255.254")},
	{MustParseAddr("0.0.0.0"), MustParseAddr("0.0.0.1"), Addr{}},
	{MustParseAddr("::"), MustParseAddr("::1"), Addr{}},
	{MustParseAddr("::%x"), MustParseAddr("::1%x"), Addr{}},
	{MustParseAddr("::1"), MustParseAddr("::2"), MustParseAddr("::")},
	{MustParseAddr("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"), Addr{}, MustParseAddr("ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffe")},
}

func TestIPNextPrev(t *testing.T) {
	doNextPrev(t)

	for _, ip := range []Addr{
		MustParseAddr("0.0.0.0"),
		MustParseAddr("::"),
	} {
		got := ip.Prev()
		if !got.isZero() {
			t.Errorf("IP(%v).Prev = %v; want zero", ip, got)
		}
	}

	var allFF [16]byte
	for i := range allFF {
		allFF[i] = 0xff
	}

	for _, ip := range []Addr{
		MustParseAddr("255.255.255.255"),
		AddrFrom16(allFF),
	} {
		got := ip.Next()
		if !got.isZero() {
			t.Errorf("IP(%v).Next = %v; want zero", ip, got)
		}
	}
}

func TestIPBitLen(t *testing.T) {
	tests := []struct {
		ip   Addr
		want int
	}{
		{Addr{}, 0},
		{MustParseAddr("0.0.0.0"), 32},
		{MustParseAddr("10.0.0.1"), 32},
		{MustParseAddr("::"), 128},
		{MustParseAddr("fed0::1"), 128},
		{MustParseAddr("::ffff:10.0.0.1"), 128},
	}
	for _, tt := range tests {
		got := tt.ip.BitLen()
		if got != tt.want {
			t.Errorf("BitLen(%v) = %d; want %d", tt.ip, got, tt.want)
		}
	}
}

func TestPrefixContains(t *testing.T) {
	tests := []struct {
		ipp  Prefix
		ip   Addr
		want bool
	}{
		{MustParsePrefix("9.8.7.6/0"), MustParseAddr("9.8.7.6"), true},
		{MustParsePrefix("9.8.7.6/16"), MustParseAddr("9.8.7.6"), true},
		{MustParsePrefix("9.8.7.6/16"), MustParseAddr("9.8.6.4"), true},
		{MustParsePrefix("9.8.7.6/16"), MustParseAddr("9.9.7.6"), false},
		{MustParsePrefix("9.8.7.6/32"), MustParseAddr("9.8.7.6"), true},
		{MustParsePrefix("9.8.7.6/32"), MustParseAddr("9.8.7.7"), false},
		{MustParsePrefix("9.8.7.6/32"), MustParseAddr("9.8.7.7"), false},
		{MustParsePrefix("::1/0"), MustParseAddr("::1"), true},
		{MustParsePrefix("::1/0"), MustParseAddr("::2"), true},
		{MustParsePrefix("::1/127"), MustParseAddr("::1"), true},
		{MustParsePrefix("::1/127"), MustParseAddr("::2"), false},
		{MustParsePrefix("::1/128"), MustParseAddr("::1"), true},
		{MustParsePrefix("::1/127"), MustParseAddr("::2"), false},
		// Zones ignored: https://go.dev/issue/51899
		{Prefix{MustParseAddr("1.2.3.4").WithZone("a"), 32}, MustParseAddr("1.2.3.4"), true},
		{Prefix{MustParseAddr("::1").WithZone("a"), 128}, MustParseAddr("::1"), true},
		// invalid IP
		{MustParsePrefix("::1/0"), Addr{}, false},
		{MustParsePrefix("1.2.3.4/0"), Addr{}, false},
		// invalid Prefix
		{PrefixFrom(MustParseAddr("::1"), 129), MustParseAddr("::1"), false},
		{PrefixFrom(MustParseAddr("1.2.3.4"), 33), MustParseAddr("1.2.3.4"), false},
		{PrefixFrom(Addr{}, 0), MustParseAddr("1.2.3.4"), false},
		{PrefixFrom(Addr{}, 32), MustParseAddr("1.2.3.4"), false},
		{PrefixFrom(Addr{}, 128), MustParseAddr("::1"), false},
		// wrong IP family
		{MustParsePrefix("::1/0"), MustParseAddr("1.2.3.4"), false},
		{MustParsePrefix("1.2.3.4/0"), MustParseAddr("::1"), false},
	}
	for _, tt := range tests {
		got := tt.ipp.Contains(tt.ip)
		if got != tt.want {
			t.Errorf("(%v).Contains(%v) = %v want %v", tt.ipp, tt.ip, got, tt.want)
		}
	}
}

func TestIPv6Accessor(t *testing.T) {
	var a [16]byte
	for i := range a {
		a[i] = uint8(i) + 1
	}
	ip := AddrFrom16(a)
	for i := range a {
		if got, want := ip.v6(uint8(i)), uint8(i)+1; got != want {
			t.Errorf("v6(%v) = %v; want %v", i, got, want)
		}
	}
}

func make128(hi, lo uint64) uint128 {
	return uint128{hi, lo}
}

func makeAddr(u uint128, bitlen uint8) Addr {
	return Addr{addr: u, bitlen: bitlen}
}

func makeAddrZone(u uint128, bitlen uint8, scopeID uint32) Addr {
	return Addr{addr: u, bitlen: bitlen, scopeID: scopeID}
}

func doNextPrev(t testing.TB) {
	for _, tt := range nextPrevTests {
		gnext, gprev := tt.ip.Next(), tt.ip.Prev()
		if gnext != tt.next {
			t.Errorf("IP(%v).Next = %v; want %v", tt.ip, gnext, tt.next)
		}
		if gprev != tt.prev {
			t.Errorf("IP(%v).Prev = %v; want %v", tt.ip, gprev, tt.prev)
		}
		if !tt.ip.Next().isZero() && tt.ip.Next().Prev() != tt.ip {
			t.Errorf("IP(%v).Next.Prev = %v; want %v", tt.ip, tt.ip.Next().Prev(), tt.ip)
		}
		if !tt.ip.Prev().isZero() && tt.ip.Prev().Next() != tt.ip {
			t.Errorf("IP(%v).Prev.Next = %v; want %v", tt.ip, tt.ip.Prev().Next(), tt.ip)
		}
	}
}
