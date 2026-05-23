// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netip

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"solod.dev/so/strings"
)

func TestPrefix(t *testing.T) {
	tests := []struct {
		prefix      string
		ip          Addr
		bits        int
		str         string
		contains    []Addr
		notContains []Addr
	}{
		{
			prefix:      "192.168.0.0/24",
			ip:          MustParseAddr("192.168.0.0"),
			bits:        24,
			contains:    mustIPs("192.168.0.1", "192.168.0.55"),
			notContains: mustIPs("192.168.1.1", "1.1.1.1"),
		},
		{
			prefix:      "192.168.1.1/32",
			ip:          MustParseAddr("192.168.1.1"),
			bits:        32,
			contains:    mustIPs("192.168.1.1"),
			notContains: mustIPs("192.168.1.2"),
		},
		{
			prefix:      "100.64.0.0/10", // CGNAT range; prefix not multiple of 8
			ip:          MustParseAddr("100.64.0.0"),
			bits:        10,
			contains:    mustIPs("100.64.0.0", "100.64.0.1", "100.81.251.94", "100.100.100.100", "100.127.255.254", "100.127.255.255"),
			notContains: mustIPs("100.63.255.255", "100.128.0.0"),
		},
		{
			prefix:      "2001:db8::/96",
			ip:          MustParseAddr("2001:db8::"),
			bits:        96,
			contains:    mustIPs("2001:db8::aaaa:bbbb", "2001:db8::1"),
			notContains: mustIPs("2001:db8::1:aaaa:bbbb", "2001:db9::"),
		},
		{
			prefix:      "0.0.0.0/0",
			ip:          MustParseAddr("0.0.0.0"),
			bits:        0,
			contains:    mustIPs("192.168.0.1", "1.1.1.1"),
			notContains: append(mustIPs("2001:db8::1"), Addr{}),
		},
		{
			prefix:      "::/0",
			ip:          MustParseAddr("::"),
			bits:        0,
			contains:    mustIPs("::1", "2001:db8::1"),
			notContains: mustIPs("192.0.2.1"),
		},
		{
			prefix:      "2000::/3",
			ip:          MustParseAddr("2000::"),
			bits:        3,
			contains:    mustIPs("2001:db8::1"),
			notContains: mustIPs("fe80::1"),
		},
	}
	for _, test := range tests {
		t.Run(test.prefix, func(t *testing.T) {
			var buf [64]byte

			prefix, err := ParsePrefix(test.prefix)
			if err != nil {
				t.Fatal(err)
			}
			if prefix.Addr() != test.ip {
				gots := prefix.Addr().String(buf[:])
				wants := test.ip.String(buf[:])
				t.Errorf("IP=%s, want %s", gots, wants)
			}
			if prefix.Bits() != test.bits {
				t.Errorf("bits=%d, want %d", prefix.Bits(), test.bits)
			}
			for _, ip := range test.contains {
				if !prefix.Contains(ip) {
					ips := ip.String(buf[:])
					t.Errorf("does not contain %s", ips)
				}
			}
			for _, ip := range test.notContains {
				if prefix.Contains(ip) {
					ips := ip.String(buf[:])
					t.Errorf("contains %s", ips)
				}
			}
			want := test.str
			if want == "" {
				want = test.prefix
			}
			if got := prefix.String(buf[:]); got != want {
				t.Errorf("prefix.String()=%q, want %q", got, want)
			}
		})
	}
}

func TestPrefixFromInvalidBits(t *testing.T) {
	v4 := MustParseAddr("1.2.3.4")
	v6 := MustParseAddr("66::66")
	tests := []struct {
		ip       Addr
		in, want int
	}{
		{v4, 0, 0},
		{v6, 0, 0},
		{v4, 1, 1},
		{v4, 33, -1},
		{v6, 33, 33},
		{v6, 127, 127},
		{v6, 128, 128},
		{v4, 254, -1},
		{v4, 255, -1},
		{v4, -1, -1},
		{v6, -1, -1},
		{v4, -5, -1},
		{v6, -5, -1},
	}
	for _, tt := range tests {
		p := PrefixFrom(tt.ip, tt.in)
		if got := p.Bits(); got != tt.want {
			t.Errorf("for (%v, %v), Bits out = %v; want %v", tt.ip, tt.in, got, tt.want)
		}
	}
}

func TestPrefixMasked(t *testing.T) {
	tests := []struct {
		prefix Prefix
		masked Prefix
	}{
		{
			prefix: MustParsePrefix("192.168.0.255/24"),
			masked: MustParsePrefix("192.168.0.0/24"),
		},
		{
			prefix: MustParsePrefix("2100::/3"),
			masked: MustParsePrefix("2000::/3"),
		},
		{
			prefix: PrefixFrom(MustParseAddr("2000::"), 129),
			masked: Prefix{},
		},
		{
			prefix: PrefixFrom(MustParseAddr("1.2.3.4"), 33),
			masked: Prefix{},
		},
	}
	for _, test := range tests {
		var buf [64]byte
		t.Run(test.prefix.String(buf[:]), func(t *testing.T) {
			got := test.prefix.Masked()
			if got != test.masked {
				gots := got.String(buf[:])
				wants := test.masked.String(buf[:])
				t.Errorf("Masked=%s, want %s", gots, wants)
			}
		})
	}
}

func TestParsePrefixError(t *testing.T) {
	tests := []struct {
		prefix string
		errstr string
	}{
		{
			prefix: "192.168.0.0",
			errstr: "no '/'",
		},
		{
			prefix: "1.257.1.1/24",
			errstr: "value >255",
		},
		{
			prefix: "1.1.1.0/q",
			errstr: "bad bits",
		},
		{
			prefix: "1.1.1.0/-1",
			errstr: "bad bits",
		},
		{
			prefix: "1.1.1.0/33",
			errstr: "out of range",
		},
		{
			prefix: "2001::/129",
			errstr: "out of range",
		},
		// Zones are not allowed: https://go.dev/issue/51899
		{
			prefix: "1.1.1.0%a/24",
			errstr: "unexpected character",
		},
		{
			prefix: "2001:db8::%a/32",
			errstr: "zones cannot be present",
		},
		{
			prefix: "1.1.1.0/+32",
			errstr: "bad bits",
		},
		{
			prefix: "1.1.1.0/-32",
			errstr: "bad bits",
		},
		{
			prefix: "1.1.1.0/032",
			errstr: "bad bits",
		},
		{
			prefix: "1.1.1.0/0032",
			errstr: "bad bits",
		},
	}
	for _, test := range tests {
		t.Run(test.prefix, func(t *testing.T) {
			_, err := ParsePrefix(test.prefix)
			if err == nil {
				t.Fatal("no error")
			}
		})
	}
}

func TestPrefixIsSingleIP(t *testing.T) {
	tests := []struct {
		ipp  Prefix
		want bool
	}{
		{ipp: MustParsePrefix("127.0.0.1/32"), want: true},
		{ipp: MustParsePrefix("127.0.0.1/31"), want: false},
		{ipp: MustParsePrefix("127.0.0.1/0"), want: false},
		{ipp: MustParsePrefix("::1/128"), want: true},
		{ipp: MustParsePrefix("::1/127"), want: false},
		{ipp: MustParsePrefix("::1/0"), want: false},
		{ipp: Prefix{}, want: false},
	}
	for _, tt := range tests {
		got := tt.ipp.IsSingleIP()
		if got != tt.want {
			t.Errorf("IsSingleIP(%v) = %v want %v", tt.ipp, got, tt.want)
		}
	}
}

func TestPrefixCompare(t *testing.T) {
	tests := []struct {
		a, b Prefix
		want int
	}{
		{Prefix{}, Prefix{}, 0},
		{Prefix{}, MustParsePrefix("1.2.3.0/24"), -1},

		{MustParsePrefix("1.2.3.0/24"), MustParsePrefix("1.2.3.0/24"), 0},
		{MustParsePrefix("fe80::/64"), MustParsePrefix("fe80::/64"), 0},

		{MustParsePrefix("1.2.3.0/24"), MustParsePrefix("1.2.4.0/24"), -1},
		{MustParsePrefix("fe80::/64"), MustParsePrefix("fe90::/64"), -1},

		{MustParsePrefix("1.2.0.0/16"), MustParsePrefix("1.2.0.0/24"), -1},
		{MustParsePrefix("fe80::/48"), MustParsePrefix("fe80::/64"), -1},

		{MustParsePrefix("1.2.3.0/24"), MustParsePrefix("fe80::/8"), -1},

		{MustParsePrefix("1.2.3.0/24"), MustParsePrefix("1.2.3.4/24"), -1},
		{MustParsePrefix("1.2.3.0/24"), MustParsePrefix("1.2.3.0/28"), -1},
	}
	for _, tt := range tests {
		var buf [64]byte
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
				t.Errorf("Compare(%q, %q) was correctly %v, but Compare(%q, %q) was %v", a, b, got, b, a, got2)
			}
		}
	}

	// And just sort.
	values := []Prefix{
		MustParsePrefix("1.2.3.0/24"), // 0
		MustParsePrefix("fe90::/64"),  // 1
		MustParsePrefix("fe80::/64"),  // 2
		MustParsePrefix("1.2.0.0/16"), // 3
		Prefix{},                      // 4
		MustParsePrefix("fe80::/48"),  // 5
		MustParsePrefix("1.2.0.0/24"), // 6
		MustParsePrefix("1.2.3.4/24"), // 7
		MustParsePrefix("1.2.3.0/28"), // 8
	}
	sorted := slices.Clone(values)
	slices.SortFunc(sorted, Prefix.Compare)
	want := []int{4, 3, 6, 0, 7, 8, 5, 2, 1} // indices of values in sorted order
	for i, v := range sorted {
		if v != values[want[i]] {
			t.Errorf("unexpected sort at index %d: got %v, want %v", i, v, values[want[i]])
		}
	}

	// Lists from
	// https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml and
	// https://www.iana.org/assignments/ipv6-address-space/ipv6-address-space.xhtml,
	// to verify that the sort order matches IANA's conventional
	// ordering.
	values = []Prefix{
		MustParsePrefix("0.0.0.0/8"),
		MustParsePrefix("127.0.0.0/8"),
		MustParsePrefix("10.0.0.0/8"),
		MustParsePrefix("203.0.113.0/24"),
		MustParsePrefix("169.254.0.0/16"),
		MustParsePrefix("192.0.0.0/24"),
		MustParsePrefix("240.0.0.0/4"),
		MustParsePrefix("192.0.2.0/24"),
		MustParsePrefix("192.0.0.170/32"),
		MustParsePrefix("198.18.0.0/15"),
		MustParsePrefix("192.0.0.8/32"),
		MustParsePrefix("0.0.0.0/32"),
		MustParsePrefix("192.0.0.9/32"),
		MustParsePrefix("198.51.100.0/24"),
		MustParsePrefix("192.168.0.0/16"),
		MustParsePrefix("192.0.0.10/32"),
		MustParsePrefix("192.175.48.0/24"),
		MustParsePrefix("192.52.193.0/24"),
		MustParsePrefix("100.64.0.0/10"),
		MustParsePrefix("255.255.255.255/32"),
		MustParsePrefix("192.31.196.0/24"),
		MustParsePrefix("172.16.0.0/12"),
		MustParsePrefix("192.0.0.0/29"),
		MustParsePrefix("192.88.99.0/24"),
		MustParsePrefix("fec0::/10"),
		MustParsePrefix("6000::/3"),
		MustParsePrefix("fe00::/9"),
		MustParsePrefix("8000::/3"),
		MustParsePrefix("0000::/8"),
		MustParsePrefix("0400::/6"),
		MustParsePrefix("f800::/6"),
		MustParsePrefix("e000::/4"),
		MustParsePrefix("ff00::/8"),
		MustParsePrefix("a000::/3"),
		MustParsePrefix("fc00::/7"),
		MustParsePrefix("1000::/4"),
		MustParsePrefix("0800::/5"),
		MustParsePrefix("4000::/3"),
		MustParsePrefix("0100::/8"),
		MustParsePrefix("c000::/3"),
		MustParsePrefix("fe80::/10"),
		MustParsePrefix("0200::/7"),
		MustParsePrefix("f000::/5"),
		MustParsePrefix("2000::/3"),
	}
	slices.SortFunc(values, func(a, b Prefix) int { return a.Compare(b) })
	var got strings.Builder
	got.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			got.WriteByte(' ')
		}
		var buf [64]byte
		got.WriteString(v.String(buf[:]))
	}
	got.WriteByte(']')
	gots := got.String()
	wants := `[0.0.0.0/8 0.0.0.0/32 10.0.0.0/8 100.64.0.0/10 127.0.0.0/8 169.254.0.0/16 172.16.0.0/12 192.0.0.0/24 192.0.0.0/29 192.0.0.8/32 192.0.0.9/32 192.0.0.10/32 192.0.0.170/32 192.0.2.0/24 192.31.196.0/24 192.52.193.0/24 192.88.99.0/24 192.168.0.0/16 192.175.48.0/24 198.18.0.0/15 198.51.100.0/24 203.0.113.0/24 240.0.0.0/4 255.255.255.255/32 ::/8 100::/8 200::/7 400::/6 800::/5 1000::/4 2000::/3 4000::/3 6000::/3 8000::/3 a000::/3 c000::/3 e000::/4 f000::/5 f800::/6 fc00::/7 fe00::/9 fe80::/10 fec0::/10 ff00::/8]`
	if gots != wants {
		t.Errorf("unexpected sort\n got: %s\nwant: %s\n", gots, wants)
	}
}

func TestPrefixMasking(t *testing.T) {
	type subtest struct {
		ip   Addr
		bits uint8
		p    Prefix
		ok   bool
	}

	// makeIPv6 produces a set of IPv6 subtests with an optional zone identifier.
	makeIPv6 := func(zone string) []subtest {
		if zone != "" {
			zone = "%" + zone
		}

		return []subtest{
			{
				ip:   MustParseAddr(fmt.Sprintf("2001:db8::1%s", zone)),
				bits: 255,
			},
			{
				ip:   MustParseAddr(fmt.Sprintf("2001:db8::1%s", zone)),
				bits: 32,
				p:    MustParsePrefix("2001:db8::/32"),
				ok:   true,
			},
			{
				ip:   MustParseAddr(fmt.Sprintf("fe80::dead:beef:dead:beef%s", zone)),
				bits: 96,
				p:    MustParsePrefix("fe80::dead:beef:0:0/96"),
				ok:   true,
			},
			{
				ip:   MustParseAddr(fmt.Sprintf("aaaa::%s", zone)),
				bits: 4,
				p:    MustParsePrefix("a000::/4"),
				ok:   true,
			},
			{
				ip:   MustParseAddr(fmt.Sprintf("::%s", zone)),
				bits: 63,
				p:    MustParsePrefix("::/63"),
				ok:   true,
			},
		}
	}

	tests := []struct {
		family   string
		subtests []subtest
	}{
		{
			family: "nil",
			subtests: []subtest{
				{
					bits: 255,
					ok:   true,
				},
				{
					bits: 16,
					ok:   true,
				},
			},
		},
		{
			family: "IPv4",
			subtests: []subtest{
				{
					ip:   MustParseAddr("192.0.2.0"),
					bits: 255,
				},
				{
					ip:   MustParseAddr("192.0.2.0"),
					bits: 16,
					p:    MustParsePrefix("192.0.0.0/16"),
					ok:   true,
				},
				{
					ip:   MustParseAddr("255.255.255.255"),
					bits: 20,
					p:    MustParsePrefix("255.255.240.0/20"),
					ok:   true,
				},
				{
					// Partially masking one byte that contains both
					// 1s and 0s on either side of the mask limit.
					ip:   MustParseAddr("100.98.156.66"),
					bits: 10,
					p:    MustParsePrefix("100.64.0.0/10"),
					ok:   true,
				},
			},
		},
		{
			family:   "IPv6",
			subtests: makeIPv6(""),
		},
		{
			family:   "IPv6 zone",
			subtests: makeIPv6("eth0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.family, func(t *testing.T) {
			for _, st := range tt.subtests {
				var buf [64]byte
				t.Run(st.p.String(buf[:]), func(t *testing.T) {
					// Ensure st.ip is not mutated.
					orig := st.ip.String(buf[:])

					p, err := st.ip.Prefix(int(st.bits))
					if st.ok && err != nil {
						t.Fatalf("failed to produce prefix: %v", err)
					}
					if !st.ok && err == nil {
						t.Fatal("expected an error, but none occurred")
					}
					if err != nil {
						t.Logf("err: %v", err)
						return
					}

					ps := p.String(buf[:])
					stps := st.p.String(buf[:])
					if !reflect.DeepEqual(p, st.p) {
						t.Errorf("prefix = %q, want %q", ps, stps)
					}

					if got := st.ip.String(buf[:]); got != orig {
						t.Errorf("IP was mutated: %q, want %q", got, orig)
					}
				})
			}
		})
	}
}

func TestPrefixOverlaps(t *testing.T) {
	var a16 [16]byte
	pfx := MustParsePrefix
	tests := []struct {
		a, b Prefix
		want bool
	}{
		{Prefix{}, pfx("1.2.0.0/16"), false},    // first zero
		{pfx("1.2.0.0/16"), Prefix{}, false},    // second zero
		{pfx("::0/3"), pfx("0.0.0.0/3"), false}, // different families

		{pfx("1.2.0.0/16"), pfx("1.2.0.0/16"), true}, // equal

		{pfx("1.2.0.0/16"), pfx("1.2.3.0/24"), true},
		{pfx("1.2.3.0/24"), pfx("1.2.0.0/16"), true},

		{pfx("1.2.0.0/16"), pfx("1.2.3.0/32"), true},
		{pfx("1.2.3.0/32"), pfx("1.2.0.0/16"), true},

		// Match /0 either order
		{pfx("1.2.3.0/32"), pfx("0.0.0.0/0"), true},
		{pfx("0.0.0.0/0"), pfx("1.2.3.0/32"), true},

		{pfx("1.2.3.0/32"), pfx("5.5.5.5/0"), true}, // normalization not required; /0 means true

		// IPv6 overlapping
		{pfx("5::1/128"), pfx("5::0/8"), true},
		{pfx("5::0/8"), pfx("5::1/128"), true},

		// IPv6 not overlapping
		{pfx("1::1/128"), pfx("2::2/128"), false},
		{pfx("0100::0/8"), pfx("::1/128"), false},

		// IPv4-mapped IPv6 addresses should not overlap with IPv4.
		{PrefixFrom(AddrFrom16(MustParseAddr("1.2.0.0").As16(a16)), 16), pfx("1.2.3.0/24"), false},

		// Invalid prefixes
		{PrefixFrom(MustParseAddr("1.2.3.4"), 33), pfx("1.2.3.0/24"), false},
		{PrefixFrom(MustParseAddr("2000::"), 129), pfx("2000::/64"), false},
	}
	for i, tt := range tests {
		if got := tt.a.Overlaps(tt.b); got != tt.want {
			t.Errorf("%d. (%v).Overlaps(%v) = %v; want %v", i, tt.a, tt.b, got, tt.want)
		}
		// Overlaps is commutative
		if got := tt.b.Overlaps(tt.a); got != tt.want {
			t.Errorf("%d. (%v).Overlaps(%v) = %v; want %v", i, tt.b, tt.a, got, tt.want)
		}
	}
}

func TestPrefixString(t *testing.T) {
	tests := []struct {
		ipp  Prefix
		want string
	}{
		{Prefix{}, "invalid Prefix"},
		{PrefixFrom(Addr{}, 8), "invalid Prefix"},
		{PrefixFrom(MustParseAddr("1.2.3.4"), 88), "invalid Prefix"},
	}

	for _, tt := range tests {
		var buf [64]byte
		if got := tt.ipp.String(buf[:]); got != tt.want {
			t.Errorf("(%#v).String() = %q want %q", tt.ipp, got, tt.want)
		}
	}
}

func mustIPs(strs ...string) []Addr {
	var res []Addr
	for _, s := range strs {
		res = append(res, MustParseAddr(s))
	}
	return res
}
