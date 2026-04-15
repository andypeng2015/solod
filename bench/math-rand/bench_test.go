package main

import (
	"math/rand/v2"
	"testing"
)

func testRand() *rand.Rand {
	return rand.New(rand.NewPCG(1, 2))
}

func Benchmark_SourceUint64(b *testing.B) {
	s := rand.NewPCG(1, 2)
	var t uint64
	for b.Loop() {
		t += s.Uint64()
	}
	sink = uint64(t)
}

func Benchmark_GlobalUint64(b *testing.B) {
	var t uint64
	for b.Loop() {
		t += rand.Uint64()
	}
	sink = t
}

func Benchmark_Uint64(b *testing.B) {
	r := testRand()
	var t uint64
	for b.Loop() {
		t += r.Uint64()
	}
	sink = t
}

func Benchmark_Int64N1e9(b *testing.B) {
	r := testRand()
	var t int64
	for b.Loop() {
		t += r.Int64N(1e9)
	}
	sink = uint64(t)
}

func Benchmark_Int64N1e18(b *testing.B) {
	r := testRand()
	var t int64
	for b.Loop() {
		t += r.Int64N(1e18)
	}
	sink = uint64(t)
}

func Benchmark_Int64N4e18(b *testing.B) {
	r := testRand()
	var t int64
	for b.Loop() {
		t += r.Int64N(4e18)
	}
	sink = uint64(t)
}

func Benchmark_Float64(b *testing.B) {
	r := testRand()
	var t float64
	for b.Loop() {
		t += r.Float64()
	}
	sink = uint64(t)
}
