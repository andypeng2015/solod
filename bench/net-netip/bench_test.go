package main

import (
	"testing"

	"solod.dev/so/net/netip"
)

func Benchmark_Parse_v4(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v4)
	}
}

func Benchmark_Parse_v6(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v6)
	}
}

func Benchmark_Parse_v6e(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v6e)
	}
}

func Benchmark_Parse_v6_v4(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v6_v4)
	}
}

func Benchmark_Parse_v6_zone(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkIP, _ = netip.ParseAddr(v6_zone)
	}
}

func Benchmark_String_v4(b *testing.B) {
	ip := netip.MustParseAddr(v4)
	buf := make([]byte, netip.MaxAddrLen)
	b.ReportAllocs()
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}

func Benchmark_String_v6(b *testing.B) {
	ip := netip.MustParseAddr(v6)
	buf := make([]byte, netip.MaxAddrLen)
	b.ReportAllocs()
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}

func Benchmark_String_v6e(b *testing.B) {
	ip := netip.MustParseAddr(v6e)
	buf := make([]byte, netip.MaxAddrLen)
	b.ReportAllocs()
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}

func Benchmark_String_v6_v4(b *testing.B) {
	ip := netip.MustParseAddr(v6_v4)
	buf := make([]byte, netip.MaxAddrLen)
	b.ReportAllocs()
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}

func Benchmark_String_v6_zone(b *testing.B) {
	ip := netip.MustParseAddr(v6_zone)
	buf := make([]byte, netip.MaxAddrLen)
	b.ReportAllocs()
	for b.Loop() {
		sinkStr = ip.String(buf)
	}
}
