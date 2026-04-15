# math/rand benchmarks

Run the benchmark:

```text
make bench name=math-rand
```

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/math-rand
cpu: Apple M1
Benchmark_SourceUint64-8    246180568     4.651 ns/op
Benchmark_GlobalUint64-8    247823476     4.847 ns/op
Benchmark_Uint64-8          268016251     4.485 ns/op
Benchmark_Int64N1e9-8       260184130     4.605 ns/op
Benchmark_Int64N1e18-8      214864120     5.596 ns/op
Benchmark_Int64N4e18-8      131153763     9.148 ns/op
Benchmark_Float64-8         271671560     4.418 ns/op
```

So:

```text
Benchmark_SourceUint64      390962254     2.836 ns/op
Benchmark_GlobalUint64      135876642     8.835 ns/op
Benchmark_Uint64            136649979     8.780 ns/op
Benchmark_Int64N1e9         131540148     9.116 ns/op
Benchmark_Int64N1e18        124685299     9.623 ns/op
Benchmark_Int64N4e18         99717466    12.07 ns/op
Benchmark_Float64           127862793     9.270 ns/op
```
