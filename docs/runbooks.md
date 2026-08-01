# CloudStoreX SRE Operations & Runbook Directory

## 1. Document Index

The CloudStoreX Site Reliability Engineering (SRE) documentation suite is organized as follows:

1. **[Observability & Telemetry Architecture](observability.md)**: OpenTelemetry distributed tracing, RED & USE methodologies, and zero-overhead instrumentation.
2. **[Monitoring & Metrics Catalogue](monitoring.md)**: Four Golden Signals and low-cardinality Prometheus metrics reference.
3. **[Alerting Guide](alerting.md)**: AlertManager routing, severity hierarchy, and silence procedures.
4. **[Performance & Load Testing](performance.md)**: Go micro-benchmarks and k6 distributed load testing profiles.
5. **[Go Runtime Profiling](profiling.md)**: Secure production `pprof` usage for CPU, Heap, Goroutine, and Mutex diagnostics.
6. **[Incident Response Runbooks](sre/runbooks.md)**: Actionable Standard Operating Procedures (SOPs) for all AlertManager alerts.

---

## 2. Escalation Matrix

| Component | First Responder (0-15 min) | Escalation Tier 2 (15-30 min) | Escalation Tier 3 (30+ min) |
| :--- | :--- | :--- | :--- |
| **API Gateway / Routing** | Backend SRE On-Call | Principal Platform Architect | Cloud Network Engineering |
| **Storage Providers** | Storage Cloud SRE | Storage Vendor Support (AWS/MinIO) | Platform Architect |
| **PostgreSQL Catalog** | Database Reliability Engineer (DBA) | Cloud Database Support | Lead Platform SRE |
| **Redis & Job Queue** | Backend SRE On-Call | Platform Infrastructure Lead | Cloud Cache Support |
