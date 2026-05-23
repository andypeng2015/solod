package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/net/netip"
	"solod.dev/so/testing"
)

//so:volatile
var sinkIP netip.Addr

//so:volatile
var sinkStr string

const (
	v4      = "192.168.1.1"
	v6      = "fd7a:115c:a1e0:ab12:4843:cd96:626b:430b"
	v6e     = "fd7a:115c::626b:430b"
	v6_v4   = "::ffff:192.168.140.255"
	v6_zone = "1:2::ffff:192.168.140.255%eth1"
)

func Parse_v4(b *testing.B) {
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v4)
	}
}

func Parse_v6(b *testing.B) {
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v6)
	}
}

func Parse_v6e(b *testing.B) {
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v6e)
	}
}

func Parse_v6_v4(b *testing.B) {
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v6_v4)
	}
}

func Parse_v6_zone(b *testing.B) {
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v6_zone)
	}
}

func String_v4(b *testing.B) {
	ip := netip.MustParseAddr(v4)
	buf := make([]byte, netip.MaxAddrLen)
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}

func String_v6(b *testing.B) {
	ip := netip.MustParseAddr(v6)
	buf := make([]byte, netip.MaxAddrLen)
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}

func String_v6e(b *testing.B) {
	ip := netip.MustParseAddr(v6e)
	buf := make([]byte, netip.MaxAddrLen)
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}

func String_v6_v4(b *testing.B) {
	ip := netip.MustParseAddr(v6_v4)
	buf := make([]byte, netip.MaxAddrLen)
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}

func String_v6_zone(b *testing.B) {
	ip := netip.MustParseAddr(v6_zone)
	buf := make([]byte, netip.MaxAddrLen)
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "Parse_v4", F: Parse_v4},
		{Name: "Parse_v6", F: Parse_v6},
		{Name: "Parse_v6e", F: Parse_v6e},
		{Name: "Parse_v6_v4", F: Parse_v6_v4},
		{Name: "Parse_v6_zone", F: Parse_v6_zone},
		{Name: "String_v4", F: String_v4},
		{Name: "String_v6", F: String_v6},
		{Name: "String_v6e", F: String_v6e},
		{Name: "String_v6_v4", F: String_v6_v4},
		{Name: "String_v6_zone", F: String_v6_zone},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
