// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package netip defines an IP address type that's a small value type.
// Building on that [Addr] type, the package also defines [AddrPort] (an
// IP address and a port) and [Prefix] (an IP address and a bit length
// prefix).
package netip

import "solod.dev/so/errors"

var ErrIP = errors.New("netip: invalid IP address")
var ErrIPv4 = errors.New("netip: invalid IPv4 address")
var ErrIPv6 = errors.New("netip: invalid IPv6 address")
var ErrIPPort = errors.New("netip: invalid ip-port")
var ErrPort = errors.New("netip: invalid port")
var ErrPrefix = errors.New("netip: invalid prefix")

var ErrNegativePrefix = errors.New("netip: negative Prefix bits")
var ErrLargePrefix = errors.New("netip: prefix length too large for IP type")
