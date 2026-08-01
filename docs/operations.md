# CloudStoreX Day-2 Operations & Incident Response Guide

This document defines operational runbooks, log querying, telemetry, and incident response procedures for CloudStoreX.

---

## 1. Structured Logging & Telemetry
CloudStoreX emits JSON logs to stdout/stderr, collected by Kubernetes logging agents and routed to CloudWatch Log Groups (`/aws/cloudstorex/prod/backend`):
- **Request ID Tracing**: Every HTTP request is tagged with an `X-Request-ID` header and logged across background jobs and database queries.
- **Health Probes**: Liveness (`/healthz`) and Readiness (`/readyz`) probe endpoints report dependency connectivity states.

---

## 2. Incident Runbooks

### Runbook A: Readiness Probe Failure (`503 SERVICE_UNAVAILABLE`)
- **Symptom**: Kubernetes removes backend pods from endpoints; user requests fail with 503 errors.
- **Diagnostic Steps**:
  1. Inspect pod logs:
     ```bash
     kubectl logs -l app.kubernetes.io/component=backend -n cloudstorex --tail=100
     ```
  2. Test readiness endpoint manually:
     ```bash
     kubectl exec -it <POD_NAME> -n cloudstorex -- wget -qO- http://localhost:8080/readyz
     ```
  3. Identify which dependency is reporting `"up" != "ready"` (PostgreSQL, Redis, Provider Registry, Policy Engine, or Background Workers).

### Runbook B: High Background Worker Queue Depth
- **Symptom**: Object retention policies or soft-delete cleanups lag behind schedule.
- **Remediation**:
  1. Check Redis worker queue length:
     ```bash
     kubectl exec -it <REDIS_POD> -n cloudstorex -- redis-cli LLEN cloudstorex:jobs:default
     ```
  2. Scale up backend worker replicas via HPA or manual deployment override.
