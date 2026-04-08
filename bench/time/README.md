# time benchmarks

Run the benchmark:

```text
make bench name=time
```

## Functions and methods

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/time
cpu: Apple M1
Benchmark_Date-8           163291910     7.016 ns/op
Benchmark_ISOWeek-8        134723946     8.905 ns/op
Benchmark_Now-8             34996177    34.30 ns/op
Benchmark_Since-8           72498975    16.58 ns/op
Benchmark_UnixNano-8        35027504    34.30 ns/op
Benchmark_Until-8           72336907    16.59 ns/op
```

So:

```text
Benchmark_Date             514257796     2.177 ns/op
Benchmark_ISOWeek          571559212     2.077 ns/op
Benchmark_Now               31011758    38.57 ns/op
Benchmark_Since             49003592    24.62 ns/op
Benchmark_UnixNano          30899164    38.29 ns/op
Benchmark_Until             51431509    24.26 ns/op
```

## Parsing and formatting

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/time
cpu: Apple M1
Benchmark_Format-8          30126373    39.12 ns/op
Benchmark_FormatCustom-8    21823570    55.42 ns/op
Benchmark_Parse-8           43960874    27.27 ns/op
Benchmark_ParseCustom-8     21568800    55.38 ns/op
```

So:

```text
Benchmark_Format           271301091     4.422 ns/op
Benchmark_FormatCustom       4692283   249.7 ns/op
Benchmark_Parse            188519773     5.623 ns/op
Benchmark_ParseCustom       26346988    45.09 ns/op
```
