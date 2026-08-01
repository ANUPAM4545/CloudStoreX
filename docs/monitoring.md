# CloudStoreX Enterprise Monitoring Guide & Prometheus Metric Catalogue

## 1. Golden Signals Monitoring

CloudStoreX monitors the Four Golden Signals across all production workloads:

1. **Latency**: Time taken to service requests (P50, P95, P99 upload/download/query duration).
2. **Traffic**: Demand placed on the system (HTTP request rate, upload/download byte throughput).
3. **Errors**: Rate of requests that fail (HTTP 5xx, provider SDK failures, storage quota breaches).
4. **Saturation**: How full the platform is (storage quota % used, background worker queue depth, DB connection pools).

---

## 2. Enterprise Prometheus Metrics Catalogue

All metrics exported at `/metrics` adhere to strict low-cardinality label rules (`workspace`, `provider`, `bucket`, `operation`, `region`, `status`). High-cardinality labels such as object keys, user IDs, or timestamps are strictly prohibited.

| Metric Name | Type | Labels | Golden Signal | Description |
| :--- | :--- | :--- | :--- | :--- |
| `cloudstorex_http_requests_total` | Counter | `method`, `path`, `status` | Traffic | Total HTTP requests handled by the API Gateway. |
| `cloudstorex_http_request_duration_seconds` | Histogram | `method`, `path`, `status` | Latency | HTTP request duration distribution in seconds. |
| `cloudstorex_http_errors_total` | Counter | `method`, `path`, `status` | Errors | Total HTTP 4xx and 5xx error responses. |
| `cloudstorex_provider_requests_total` | Counter | `provider`, `operation`, `region` | Traffic | Total storage operations sent to providers (aws-s3, minio). |
| `cloudstorex_provider_errors_total` | Counter | `provider`, `operation`, `region` | Errors | Failed storage provider operations. |
| `cloudstorex_provider_latency_seconds` | Histogram | `provider`, `operation`, `region` | Latency | Latency distribution of provider operations. |
| `cloudstorex_provider_health_status` | Gauge | `provider`, `region` | Saturation / Health | Health check status of provider (1=Healthy, 0=Unreachable/Degraded). |
| `cloudstorex_policy_evaluations_total` | Counter | `workspace`, `policy_id` | Traffic | Total dynamic routing policy evaluations executed. |
| `cloudstorex_policy_evaluation_latency_seconds` | Histogram | `workspace`, `policy_id` | Latency | Time spent evaluating dynamic routing rules. |
| `cloudstorex_policy_fallback_routes_total` | Counter | `workspace`, `reason` | Errors / Routing | Number of requests routed to default provider due to policy fallback. |
| `cloudstorex_storage_upload_bytes_total` | Counter | `workspace`, `provider`, `bucket` | Traffic | Total bytes uploaded across workspaces. |
| `cloudstorex_storage_download_bytes_total` | Counter | `workspace`, `provider`, `bucket` | Traffic | Total bytes downloaded across workspaces. |
| `cloudstorex_jobs_enqueued_total` | Counter | `queue`, `job_type` | Traffic | Total background jobs submitted to Redis queue. |
| `cloudstorex_jobs_completed_total` | Counter | `queue`, `job_type` | Traffic | Total successfully completed background jobs. |
| `cloudstorex_jobs_failed_total` | Counter | `queue`, `job_type` | Errors | Total background jobs failed after retries. |
| `cloudstorex_jobs_queue_depth` | Gauge | `queue` | Saturation | Current pending jobs backlog in queue. |
| `cloudstorex_jobs_execution_duration_seconds` | Histogram | `queue`, `job_type` | Latency | Background job execution time. |
| `cloudstorex_metadata_object_catalog_size` | Gauge | `workspace`, `bucket` | Saturation | Number of tracked objects in Postgres catalog. |
| `cloudstorex_postgres_connections` | Gauge | `status` (`active`, `idle`) | Saturation | Open PostgreSQL database connection pool count. |

---

## 3. Label Cardinality Enforcement

To prevent Prometheus memory exhaustion:
- `path`: Normalized route patterns (`/api/v1/storage/buckets/:bucket/objects/*key` instead of raw URLs).
- `workspace`: Low-cardinality tenant UUIDs.
- `provider`: Explicit provider IDs (`aws-s3`, `minio`).
