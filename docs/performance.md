# CloudStoreX Performance Benchmark & SLA Methodology

This document establishes the load testing methodology, workload profiles, SLA targets, and benchmarking tools for CloudStoreX (**Refinement 11**).

---

## 1. Workload Profiles

CloudStoreX evaluates control-plane and data-plane performance across three canonical workload profiles:

### A. Metadata-Heavy Workload (10 KB Objects)
- **Characteristics**: High request frequency (QPS), small payloads, database-intensive read/write operations (tagging, listing, metadata lookups, lifecycle evaluations).
- **Primary Bottleneck**: PostgreSQL transaction latency and Redis caching throughput.
- **Target Metrics**: Request latency, cache hit ratio (>90%), DB connection pool utilization.

### B. Typical Enterprise Object Workload (10 MB Objects)
- **Characteristics**: Mixed metadata and data transfer, standard document/media file sizes, single-part presigned upload/download streams.
- **Primary Bottleneck**: Network bandwidth, S3/MinIO provider latency, Policy Engine routing speed.
- **Target Metrics**: End-to-end upload/download time, Policy Engine routing latency (<5ms).

### C. Large Multipart Workload (1 GB Objects)
- **Characteristics**: Concurrent multipart upload streams, high network IO, sustained storage throughput.
- **Primary Bottleneck**: Provider bandwidth, multipart completion latency.
- **Target Metrics**: Sustained GB/sec throughput, multipart completion error rates.

---

## 2. SLA Targets

| Workload Profile | Operation | Target p50 Latency | Target p95 Latency | Target p99 Latency |
| :--- | :--- | :--- | :--- | :--- |
| **Metadata-Heavy** | Object List / Tag / Lookup | `< 10 ms` | `< 45 ms` | `< 100 ms` |
| **Typical (10 MB)** | Presigned URL Generation | `< 5 ms` | `< 15 ms` | `< 30 ms` |
| **Typical (10 MB)** | Full Upload / Download | `< 250 ms` | `< 750 ms` | `< 1500 ms` |
| **Routing Decision** | Policy Engine Provider Selection | `< 1 ms` | `< 3 ms` | `< 10 ms` |

---

## 3. Tooling & Load Testing Scripts

### A. API Latency Benchmarking with `hey`
Test metadata lookup throughput and latency distribution:
```bash
hey -n 10000 -c 50 \
  -H "Authorization: Bearer <TEST_JWT_TOKEN>" \
  http://api.cloudstorex.local/api/v1/storage/buckets/test-bucket/objects
```

### B. Load Testing with `k6`
Example `k6` script (`tests/performance/load_test.js`) for presigned upload generation:
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '1m', target: 50 },
    { duration: '3m', target: 200 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(99)<100'], // 99% of requests must complete below 100ms
  },
};

export default function () {
  let url = 'http://api.cloudstorex.local/api/v1/storage/buckets/test-bucket/presigned-upload/perf-object.bin';
  let params = {
    headers: { 'Authorization': 'Bearer <TEST_JWT_TOKEN>' },
  };
  let res = http.post(url, null, params);
  check(res, {
    'status is 200': (r) => r.status === 200,
    'has presigned URL': (r) => r.json('data') !== undefined,
  });
  sleep(0.1);
}
```
