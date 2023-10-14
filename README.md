# memcache
Go in-memory generic key-value cache.

## benchmarks
```
goos: darwin
goarch: amd64
pkg: github.com/wafer-bw/memcache
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkCache_Set/set_to_cache_with_100_keys-6                 18182341                66.23 ns/op            8 B/op          1 allocs/op
BenchmarkCache_Set/set_to_cache_with_1000_keys-6                17511355                65.91 ns/op            8 B/op          1 allocs/op
BenchmarkCache_Set/set_to_cache_with_10000_keys-6               16158224                99.20 ns/op            8 B/op          1 allocs/op
BenchmarkCache_Set/set_to_cache_with_100000_keys-6               8444706               124.7 ns/op             8 B/op          1 allocs/op
BenchmarkCache_Get/get_from_cache_with_100_keys-6               37762674                27.23 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Get/get_from_cache_with_1000_keys-6              38833897                30.75 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Get/get_from_cache_with_10000_keys-6             31752844                37.53 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Get/get_from_cache_with_100000_keys-6            16703668                63.06 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Delete/delete_from_cache_with_100_keys-6         33022230                30.68 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Delete/delete_from_cache_with_1000_keys-6        38775711                30.56 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Delete/delete_from_cache_with_10000_keys-6       38600520                30.61 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Delete/delete_from_cache_with_100000_keys-6      39075007                31.32 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Flush/flush_cache_with_100_keys-6                46152319                25.98 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Flush/flush_cache_with_1000_keys-6               46204053                25.95 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Flush/flush_cache_with_10000_keys-6              45028154                26.08 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Flush/flush_cache_with_100000_keys-6             38557881                26.56 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/wafer-bw/memcache    22.558s
```