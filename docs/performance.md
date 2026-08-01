# CloudStoreX Performance Engineering & Tuning Guide

## 1. Go Micro-Benchmarking Suite

CloudStoreX maintains a comprehensive micro-benchmarking suite across its core engines:
- `internal/storage/benchmark_test.go`: Measures synchronous and parallel object upload throughput and memory allocations.
- `internal/policy/engine/benchmark_test.go`: Measures rule evaluation latency and provider routing resolution.
- `internal/metadata/service/benchmark_test.go`: Measures object catalog lookup and multi-tenant search performance.

### Running Go Benchmarks

Execute the benchmarking suite with memory allocation reporting:

```bash
go test -bench=. -benchmem -run=^$ ./internal/storage ./internal/policy/engine ./internal/metadata/service
```

#### Verified Baseline Benchmarks (Apple M2 / 8-core ARM64)
- **Policy Engine Resolution (`BenchmarkPolicyEngine_ResolveParallel`)**: ~1,270 ns/op (0.0012 ms), 23 allocs/op (1,368 B/op).
- **Metadata Catalog Bucket Lookup (`BenchmarkMetadataService_FindBucketByName`)**: ~509 ns/op (0.0005 ms), 2 allocs/op (176 B/op).
- **Storage Object Upload Pipeline (`BenchmarkStorageService_UploadObjectParallel`)**: ~286 µs/op (0.28 ms), 22 allocs/op (2,101 B/op).

---

## 2. Distributed k6 Load Testing Profiles

In `deploy/observability/k6/`, five standardized load testing scripts validate platform scalability:

| Script | Purpose | Virtual Users (VUs) | Duration | Golden SLO Target |
| :--- | :--- | :--- | :--- | :--- |
| `smoke.js` | Rapid sanity check | 5 VUs | 30 seconds | 0% errors, P95 < 200ms |
| `load.js` | Normal daily traffic | 50 VUs | 5 minutes | < 0.1% errors, P95 < 300ms |
| `stress.js` | Saturation & breaking point | Up to 200 VUs | 10 minutes | Graceful backpressure |
| `spike.js` | Sudden 10x traffic surge | 10 -> 200 -> 10 VUs | 3 minutes | No panic/crash, recovery < 30s |
| `soak.js` | Endurance & memory leak test | 40 VUs | 4+ hours | Flat RSS memory profile |

### Executing Load Tests

Run k6 with Prometheus Remote Write output:

```bash
k6 run --out experimental-prometheus-rw deploy/observability/k6/load.js
```
