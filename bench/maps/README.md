# maps benchmarks

Requires GCC/Clang and mimalloc (for heap allocations in So). If mimalloc isn't available, the benchmarks will use the default libc allocator, which is much slower.

Run the benchmark:

```text
make bench name=maps
```

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/maps
cpu: Apple M1

Benchmark_IntSet-8       31677    35580 ns/op     74264 B/op    20 allocs/op
Benchmark_IntGet-8      218179     5573 ns/op         0 B/op     0 allocs/op
Benchmark_IntHas-8      211342     5660 ns/op         0 B/op     0 allocs/op
Benchmark_IntDelete-8    50260    23892 ns/op     36944 B/op     5 allocs/op
Benchmark_IntSetDel-8    63715    18956 ns/op         0 B/op     0 allocs/op

Benchmark_StrSet-8       24082    48677 ns/op    108760 B/op    20 allocs/op
Benchmark_StrGet-8      134481     8990 ns/op         0 B/op     0 allocs/op
Benchmark_StrHas-8      139606    10174 ns/op         0 B/op     0 allocs/op
Benchmark_StrDelete-8    34094    33878 ns/op     54608 B/op     5 allocs/op
Benchmark_StrSetDel-8    45928    26323 ns/op         0 B/op     0 allocs/op
```

So (mimalloc):

```text
Benchmark_IntSet         20305    56696 ns/op     98112 B/op    27 allocs/op
Benchmark_IntGet        734978     1638 ns/op         0 B/op     0 allocs/op
Benchmark_IntHas        783085     1596 ns/op         0 B/op     0 allocs/op
Benchmark_IntDelete      30958    38556 ns/op     73728 B/op     6 allocs/op
Benchmark_IntSetDel     263752     4628 ns/op       192 B/op     3 allocs/op

Benchmark_StrSet         18986    71879 ns/op    130816 B/op    27 allocs/op
Benchmark_StrGet        117218    10206 ns/op         0 B/op     0 allocs/op
Benchmark_StrHas        119547    10135 ns/op         0 B/op     0 allocs/op
Benchmark_StrDelete      23670    50111 ns/op     98304 B/op     6 allocs/op
Benchmark_StrSetDel     126256     9544 ns/op       256 B/op     3 allocs/op
```

So (arena):

```text
Benchmark_IntSet         21157    57661 ns/op     98112 B/op    27 allocs/op
Benchmark_IntGet        752787     1583 ns/op         0 B/op     0 allocs/op
Benchmark_IntHas        779625     1532 ns/op         0 B/op     0 allocs/op
Benchmark_IntDelete      30339    38821 ns/op     73728 B/op     6 allocs/op
Benchmark_IntSetDel     264235     4530 ns/op       192 B/op     3 allocs/op

Benchmark_StrSet         18884    63500 ns/op    130816 B/op    27 allocs/op
Benchmark_StrGet        119041    10083 ns/op         0 B/op     0 allocs/op
Benchmark_StrHas        118683    10203 ns/op         0 B/op     0 allocs/op
Benchmark_StrDelete      24212    49507 ns/op     98304 B/op     6 allocs/op
Benchmark_StrSetDel     125724     9667 ns/op       256 B/op     3 allocs/op
```
