// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netip

import (
	"solod.dev/so/bytealg"
	"solod.dev/so/cmp"
	"solod.dev/so/strconv"
)

// Maximum length of an ip:port string.
const (
	MaxPortLen         = 5                                   // 65535
	MaxAddrPort4Len    = MaxAddr4Len + 1 + MaxPortLen        // ip:port
	MaxAddrPort4In6Len = MaxAddr4In6Len + 2 + 1 + MaxPortLen // [ip]:port
	MaxAddrPort6Len    = MaxAddr6Len + 2 + 1 + MaxPortLen    // [ip]:port
	MaxAddrPortLen     = MaxAddrPort6Len
)

// AddrPort is an IP and a port number.
type AddrPort struct {
	ip   Addr
	port uint16
}

// AddrPortFrom returns an [AddrPort] with the provided IP and port.
// It does not allocate.
func AddrPortFrom(ip Addr, port uint16) AddrPort { return AddrPort{ip: ip, port: port} }

// Addr returns p's IP address.
func (p AddrPort) Addr() Addr { return p.ip }

// Port returns p's port.
func (p AddrPort) Port() uint16 { return p.port }

type addrPortParts struct {
	ip   string
	port string
	v6   bool
}

// splitAddrPort splits s into an IP address string and a port
// string. It splits strings shaped like "foo:bar" or "[foo]:bar",
// without further validating the substrings. v6 indicates whether the
// ip string should parse as an IPv6 address or an IPv4 address, in
// order for s to be a valid ip:port string.
func splitAddrPort(s string) (addrPortParts, error) {
	i := bytealg.LastIndexByteString(s, ':')
	if i == -1 {
		return addrPortParts{}, ErrIPPort
	}

	var v6 bool
	ip, port := s[:i], s[i+1:]
	if len(ip) == 0 {
		return addrPortParts{}, ErrIP
	}
	if len(port) == 0 {
		return addrPortParts{}, ErrPort
	}
	if ip[0] == '[' {
		if len(ip) < 2 || ip[len(ip)-1] != ']' {
			return addrPortParts{}, ErrIP
		}
		ip = ip[1 : len(ip)-1]
		v6 = true
	}

	return addrPortParts{ip: ip, port: port, v6: v6}, nil
}

// ParseAddrPort parses s as an [AddrPort].
//
// It doesn't do any name resolution: both the address and the port
// must be numeric.
func ParseAddrPort(s string) (AddrPort, error) {
	var ipp AddrPort
	parts, err := splitAddrPort(s)
	if err != nil {
		return ipp, err
	}
	port16, err := strconv.ParseUint(parts.port, 10, 16)
	if err != nil {
		return ipp, ErrPort
	}
	ipp.port = uint16(port16)
	ipp.ip, err = ParseAddr(parts.ip)
	if err != nil {
		return AddrPort{}, err
	}
	if parts.v6 && ipp.ip.Is4() {
		return AddrPort{}, ErrIPPort
	} else if !parts.v6 && ipp.ip.Is6() {
		return AddrPort{}, ErrIPPort
	}
	return ipp, nil
}

// MustParseAddrPort calls [ParseAddrPort](s) and panics on error.
// It is intended for use in tests with hard-coded strings.
func MustParseAddrPort(s string) AddrPort {
	ip, err := ParseAddrPort(s)
	if err != nil {
		panic(err)
	}
	return ip
}

// IsValid reports whether p.Addr() is valid.
// All ports are valid, including zero.
func (p AddrPort) IsValid() bool { return p.ip.IsValid() }

// Compare returns an integer comparing two AddrPorts.
// The result will be 0 if p == p2, -1 if p < p2, and +1 if p > p2.
// AddrPorts sort first by IP address, then port.
func (p AddrPort) Compare(p2 AddrPort) int {
	if c := p.Addr().Compare(p2.Addr()); c != 0 {
		return c
	}
	port1 := p.Port()
	port2 := p2.Port()
	return cmp.Compare(port1, port2)
}

// AppendText implements the [encoding.TextAppender] interface.
// Requires at least [MaxAddrPortLen] bytes of spare capacity in b.
func (p AddrPort) AppendText(b []byte) ([]byte, error) {
	return p.appendTo(b), nil
}

// String returns a string representation of p.
// buf length must be at least [MaxAddrPortLen].
func (p AddrPort) String(buf []byte) string {
	b := buf[:0]
	switch p.ip.bitlen {
	case z0:
		return "invalid AddrPort"
	case z4:
		b = p.ip.appendTo4(b)
	default:
		if p.ip.Is4In6() {
			b = append(b, '[')
			b = p.ip.appendTo4In6(b)
		} else {
			b = append(b, '[')
			b = p.ip.appendTo6(b)
		}
		b = append(b, ']')
	}
	b = append(b, ':')
	b = strconv.AppendUint(b, uint64(p.port), 10)
	return string(b)
}

// appendTo appends a text encoding of p
// to b and returns the extended buffer.
func (p AddrPort) appendTo(b []byte) []byte {
	switch p.ip.bitlen {
	case z0:
		return b
	case z4:
		b = p.ip.appendTo4(b)
	default:
		b = append(b, '[')
		if p.ip.Is4In6() {
			b = p.ip.appendTo4In6(b)
		} else {
			b = p.ip.appendTo6(b)
		}
		b = append(b, ']')
	}
	b = append(b, ':')
	b = strconv.AppendUint(b, uint64(p.port), 10)
	return b
}
