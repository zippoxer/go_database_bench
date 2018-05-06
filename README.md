# Go Database Benchmarks

This is a test suite for benchmarking various embedded databases written in Go. We don't store raw bytes, we store structured data. If a database doesn't serialize structured data, we serialize with `encoding/json` before storing.

## Databases

- [BadgerDB](https://github.com/dgraph-io/badger) - a key-value database written in Go
- [Bow](https://github.com/zippoxer/bow) - a database for structured data, powered by BadgerDB
- [Storm](https://github.com/asdine/storm) - a database for structured data, powered by [bbolt](https://github.com/coreos/bbolt)

## Results

5/6/2018 - Go 1.10.1 on a 64-bit Windows 10 machine with i7-6700K:

```
BenchmarkBadgerPut-8              100000             16203 ns/op            2830 B/op         67 allocs/op
BenchmarkBadgerGet-8              200000              5618 ns/op            1381 B/op         27 allocs/op
BenchmarkBadgerIter-8                100          25624041 ns/op         3480579 B/op      65818 allocs/op
BenchmarkBowPut-8                 100000             16563 ns/op            3055 B/op         73 allocs/op
BenchmarkBowGet-8                 200000              6209 ns/op            1590 B/op         33 allocs/op
BenchmarkBowIter-8                   100          28916476 ns/op         3729215 B/op      71552 allocs/op
BenchmarkStormPut-8                   20         116219155 ns/op           10764 B/op         86 allocs/op
BenchmarkStormGet-8               300000              6986 ns/op            1199 B/op         20 allocs/op
BenchmarkStormIter-8                 100          62039747 ns/op         4117739 B/op      65756 allocs/op
```

## Issues

I'm not sure why Storm is too slow in `BenchmarkStormPut`. I wouldn't rely on this result yet, I think it needs further investigation.