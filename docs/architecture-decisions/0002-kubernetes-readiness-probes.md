# ADR 0002: Comprehensive Kubernetes Readiness Probes Over Simple Ping Checks

## Status
Accepted (Epic 11)

## Context
In Kubernetes deployments, liveness probes determine whether a pod should be restarted, while readiness probes determine whether a pod should receive ingress traffic. Historically, simple HTTP ping endpoints (`200 OK` if the process is running) were used for both checks. However, if a backend pod loses connectivity to PostgreSQL, Redis, or has an uninitialized Policy Engine, routing user upload requests to that pod results in cascading 5xx HTTP errors.

## Decision
We decouple liveness (`/healthz`, `/livez`) from readiness (`/readyz`, `/api/v1/ready`):
- **/healthz**: Returns `200 OK` immediately if the HTTP server thread is running.
- **/readyz**: Performs an active dependency state inspection:
  1. Checks PostgreSQL database reachability (`sqlDB.Ping()`).
  2. Checks Redis cache/queue reachability (`a.RedisClient.Ping()`).
  3. Validates that the Provider Registry is initialized and non-empty.
  4. Validates that the Policy Engine is initialized.
  5. Confirms Background Worker Queue dispatcher connectivity.
- If any check fails, `/readyz` returns HTTP `503 Service Unavailable` with code `SERVICE_UNAVAILABLE`.

## Consequences
### Positive
- **Zero Traffic Loss**: Pods experiencing transient database dropouts or uninitialized provider registries are immediately removed from Kubernetes Service endpoints.
- **Accurate Rolling Deployments**: Kubernetes will not route traffic to new replicas during deployment until all dependencies and engines are ready.

### Negative / Mitigations
- **Probe Latency**: Checks require quick network pings. Timeout policies on probes (`timeoutSeconds: 3`, `periodSeconds: 10`) prevent thread blocking.
