# CloudStoreX Go Runtime Profiling & Diagnostics Guide (`pprof`)

## 1. Production Profiling Security

To prevent Denial-of-Service (DoS) and information disclosure, CloudStoreX runtime profiling endpoints are **disabled by default**.

### Enabling Profiling in Production
Set the following environment variables on the backend container:
```yaml
env:
  - name: ENABLE_PPROF
    value: "true"
  - name: PPROF_TOKEN
    valueFrom:
      secretKeyRef:
        name: cloudstorex-pprof-secret
        key: token
```

---

## 2. Capturing Runtime Profiles

When enabled, use `go tool pprof` or curl with the `X-Pprof-Token` header:

### 1. 30-Second CPU Profile
```bash
go tool pprof -http=:8081 "http://localhost:8080/debug/pprof/profile?seconds=30&token=secret-admin-token"
```
Analyzes CPU hotspots during high upload throughput or policy rule evaluation.

### 2. Heap Memory Allocation Profile
```bash
go tool pprof -http=:8081 "http://localhost:8080/debug/pprof/heap?token=secret-admin-token"
```
Identifies memory leaks, large object buffer allocations, and GC pressure.

### 3. Goroutine Leak Inspection
```bash
go tool pprof "http://localhost:8080/debug/pprof/goroutine?token=secret-admin-token"
```
Detects leaking goroutines in background job workers or stuck HTTP connections.

### 4. Mutex & Block Contention Profile
```bash
go tool pprof "http://localhost:8080/debug/pprof/mutex?token=secret-admin-token"
```
Identifies lock contention across provider registries and policy caches.
