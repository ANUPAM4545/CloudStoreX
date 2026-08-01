# CloudStoreX Enterprise Observability & Telemetry Architecture

## 1. Overview & Architecture

CloudStoreX implements a zero-overhead, cloud-native observability pipeline based on **OpenTelemetry (OTel)**, **Prometheus**, and **Grafana**. Every API request, background job execution, storage provider interaction, metadata query, and policy decision is instrumented to produce **RED (Rate, Errors, Duration)** and **USE (Utilization, Saturation, Errors)** signals.

```text
                  Users & Application Clients
                              │
                    HTTP / REST API Gateway
                              │
           ┌──────────────────┴──────────────────┐
           │   OpenTelemetry Middleware & SDK     │
           │   - Distributed W3C Trace Context   │
           │   - RED Signal Metrics Collector     │
           │   - Structured JSON Logger Injection │
           └──────────────────┬──────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│ Storage &     │     │ Policy Engine │     │ Background    │
│ Provider Tier │     │ & Routing     │     │ Job Workers   │
└───────┬───────┘     └───────┬───────┘     └───────┬───────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────────────────────────────────────────────────┐
│              OpenTelemetry / Prometheus Pipeline          │
└───────┬─────────────────────┬─────────────────────┬───────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│  Prometheus   │     │ OTel Collector│     │ Log Aggregator│
│  Metric Store │     │ / Jaeger / OT │     │ / Elasticsearch│
└───────┬───────┘     └───────┬───────┘     └───────┬───────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              ▼
                  ┌───────────────────────┐
                  │   Grafana Dashboards  │
                  │  & AlertManager Rules │
                  └───────────────────────┘
```

---

## 2. Distributed Tracing & W3C Trace Context Propagation

CloudStoreX enforces strict **W3C Trace Context (`traceparent`, `tracestate`)** propagation across all synchronous HTTP routes and asynchronous background workers.

1. **HTTP Ingress**: `middleware.OpenTelemetryMiddleware` intercepts every incoming HTTP request. If a W3C trace header is present, it joins the existing trace; otherwise, it generates a new root 128-bit `TraceID` and 64-bit `SpanID`.
2. **Context Propagation**: The OpenTelemetry trace span is attached to the `context.Context` embedded within `gin.Context`.
3. **Child Spans**: Business domain services (`StorageService`, `PolicyEngine`, `MetadataService`, and provider adapters) create child spans using `tracer.Start(ctx, "OperationName")`, linking parent-child hierarchies without coupling to storage logic.
4. **OTLP Exporter**: Traces are exported via OTLP HTTP/gRPC to Jaeger, OpenTelemetry Collector, or cloud-managed trace stores.

---

## 3. RED Methodology (Rate, Errors, Duration)

The **RED Methodology** is applied to all request-driven services (HTTP REST API, Storage Service, and Policy Engine):

- **Rate**: Measured via `cloudstorex_http_requests_total`, `cloudstorex_provider_requests_total`, and `cloudstorex_policy_evaluations_total`.
- **Errors**: Measured via `cloudstorex_http_errors_total`, `cloudstorex_provider_errors_total`, and `cloudstorex_policy_fallback_routes_total`.
- **Duration**: Measured via latency histograms (`cloudstorex_http_request_duration_seconds`, `cloudstorex_provider_latency_seconds_bucket`, and `cloudstorex_policy_evaluation_latency_seconds_bucket`).

---

## 4. USE Methodology (Utilization, Saturation, Errors)

The **USE Methodology** is applied to infrastructure resources (Kubernetes pods, Go runtime, PostgreSQL pool, Redis queues):

- **Utilization**:
  - CPU & Memory utilization (`container_memory_working_set_bytes`, `node_cpu_seconds_total`).
  - PostgreSQL active connections (`cloudstorex_postgres_connections{status="active"}`).
- **Saturation**:
  - Background worker queue depth (`cloudstorex_jobs_queue_depth`).
  - Goroutine count and thread creation (`go_goroutines`, `go_threads`).
  - PostgreSQL connection waiting threads.
- **Errors**:
  - Unreachable database or provider health failures (`cloudstorex_provider_health_status == 0`).
  - Background worker job failures (`cloudstorex_jobs_failed_total`).

---

## 5. Production Security & pprof Profiling

In compliance with enterprise security hardening (Refinement 7), Go runtime profiling endpoints (`/debug/pprof/`) are **disabled by default**.
- **Activation**: Requires setting environment variable `ENABLE_PPROF=true`.
- **Authentication**: When `PPROF_TOKEN` is set, requests must present `X-Pprof-Token` header or query parameter to access heap, goroutine, CPU, or mutex profiles.
