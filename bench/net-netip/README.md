# net/netip benchmarks

Run the benchmark:

```text
make bench name=net-netip
```

## Parse

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/net-netip
cpu: Apple M1
Benchmark_Parse_v4-8         65076201      17.94 ns/op
Benchmark_Parse_v6-8         14723881      81.31 ns/op
Benchmark_Parse_v6e-8        25891828      47.03 ns/op
Benchmark_Parse_v6_v4-8      24727015      48.32 ns/op
Benchmark_Parse_v6_zone-8    19207452      63.59 ns/op
```

So:

```text
Benchmark_Parse_v4           73002654      16.21 ns/op
Benchmark_Parse_v6           21485738      54.65 ns/op
Benchmark_Parse_v6e          37168963      32.68 ns/op
Benchmark_Parse_v6_v4        30818223      40.36 ns/op
Benchmark_Parse_v6_zone         63265    19087 ns/op
```

Parsing an IPv6 address with a zone in So involves an `if_nametoindex` syscall — that's why it's so much slower.

## String

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/net-netip
cpu: Apple M1
Benchmark_String_v4-8        60287556    19.81 ns/op    16 B/op    1 allocs/op
Benchmark_String_v6-8        22197643    53.33 ns/op    48 B/op    1 allocs/op
Benchmark_String_v6e-8       21842838    55.04 ns/op    24 B/op    1 allocs/op
Benchmark_String_v6_v4-8     52727067    22.84 ns/op    24 B/op    1 allocs/op
Benchmark_String_v6_zone-8   20261955    60.12 ns/op    24 B/op    1 allocs/op
```

So:

```text
Benchmark_String_v4         136978639     8.683 ns/op    0 B/op    0 allocs/op
Benchmark_String_v6          68630253    17.38 ns/op     0 B/op    0 allocs/op
Benchmark_String_v6e         89585665    16.05 ns/op     0 B/op    0 allocs/op
Benchmark_String_v6_v4      103168486    11.35 ns/op     0 B/op    0 allocs/op
Benchmark_String_v6_zone     84967782    14.04 ns/op     0 B/op    0 allocs/op
```
