# CloudStoreX SRE Incident Response Runbooks & Standard Operating Procedures (SOPs)

This document contains actionable runbooks for all CloudStoreX AlertManager alerts.

---

## 1. StorageProviderUnreachable (Critical)

- **Description**: A storage provider (`aws-s3` or `minio`) has failed continuous health checks for > 2 minutes (`cloudstorex_provider_health_status == 0`).
- **Symptoms**: Upload, download, and delete requests targeting the provider fail or trigger policy fallback.
- **Impact**: Object operations for buckets mapped to this provider fail; degraded data availability.
- **Diagnosis Steps**:
  1. Check provider health metrics:
     ```bash
     curl -s http://localhost:8080/metrics | grep cloudstorex_provider_health_status
     ```
  2. Verify network connectivity from CloudStoreX pods to provider endpoint:
     ```bash
     kubectl exec -it deploy/cloudstorex-backend -- curl -vI https://s3.us-east-1.amazonaws.com
     ```
  3. Inspect provider authentication logs for credentials expiry.
- **Remediation Steps**:
  1. If MinIO pod is down: `kubectl rollout restart statefulset/cloudstorex-minio`.
  2. If AWS S3 credentials expired: Update Kubernetes secret `cloudstorex-provider-creds` and restart backend pod.
  3. Temporarily route traffic to backup provider by disabling failing routing policy via API:
     ```bash
     curl -X POST http://localhost:8080/api/v1/policies/<policy-id>/disable -H "Authorization: Bearer <token>"
     ```
- **Escalation Path**: Page Platform Infrastructure On-Call -> Storage Cloud Vendor Support.

---

## 2. HTTPHighErrorRateCritical (Critical)

- **Description**: API 5xx server error rate exceeds 10% for > 5 minutes.
- **Symptoms**: Users experience HTTP 500/502/503 errors across REST API endpoints.
- **Impact**: System-wide service degradation or outage.
- **Diagnosis Steps**:
  1. Check error breakdown by path in Grafana Platform Overview dashboard.
  2. Check recent application error logs:
     ```bash
     kubectl logs -l app=cloudstorex-backend --tail=200 | jq 'select(.level=="ERROR")'
     ```
  3. Check database and Redis connectivity.
- **Remediation Steps**:
  1. If database connection pool is exhausted, scale up PostgreSQL max connections or restart backend pods.
  2. If caused by a bad deployment, execute immediate GitOps/Helm rollback:
     ```bash
     helm rollback cloudstorex 1
     ```
- **Escalation Path**: Page Backend SRE On-Call -> Principal Platform Architect.

---

## 3. PostgresDatabaseUnreachable (Critical)

- **Description**: Zero active PostgreSQL connections (`cloudstorex_postgres_connections == 0`) for > 1 minute.
- **Symptoms**: All metadata queries, user auth, and storage requests fail with 500 errors.
- **Impact**: Complete control-plane outage.
- **Diagnosis Steps**:
  1. Check PostgreSQL pod status: `kubectl get pods -l app=postgres`.
  2. Check Postgres disk usage and logs: `kubectl logs -l app=postgres --tail=100`.
- **Remediation Steps**:
  1. If Postgres pod is OOMKilled or crashed, increase memory limits in Helm chart and restart.
  2. If disk is full, expand PVC storage volume.
- **Escalation Path**: Page DBA / Database Reliability Engineer On-Call.

---

## 4. StorageUploadLatencyHigh (High)

- **Description**: P95 object upload latency exceeds 5.0 seconds for > 5 minutes.
- **Symptoms**: Large uploads time out or hang; slow application responsiveness.
- **Impact**: Poor end-user performance.
- **Diagnosis Steps**:
  1. Check if latency is isolated to a specific provider (`provider="aws-s3"` vs `"minio"`).
  2. Inspect network bandwidth saturation between K8s nodes and storage buckets.
- **Remediation Steps**:
  1. Enable multipart upload concurrency tuning for large payloads.
  2. If AWS S3 region throttling occurs, adjust routing policy to distribute across regions.
- **Escalation Path**: Page Cloud Network SRE On-Call.

---

## 5. BackgroundWorkerQueueDepthHigh (High)

- **Description**: Redis job queue backlog (`cloudstorex_jobs_queue_depth`) exceeds 100 pending jobs for > 10 minutes.
- **Symptoms**: Object lifecycle expiration, retention cleanup, and audit event processing are delayed.
- **Impact**: Delayed background processing; quota usage updates lag.
- **Diagnosis Steps**:
  1. Check worker logs for repeating errors or deadlocks:
     ```bash
     kubectl logs -l app=cloudstorex-backend --tail=100 | grep "Job execution failed"
     ```
- **Remediation Steps**:
  1. Scale up backend worker replicas:
     ```bash
     kubectl scale deploy cloudstorex-backend --replicas=5
     ```
  2. Clear dead-letter jobs if blocked by poisoned messages.
- **Escalation Path**: Page Backend SRE On-Call.

---

## 6. ProviderErrorRateHigh (High)

- **Description**: Storage provider error rate > 5% for > 5 minutes.
- **Symptoms**: Frequent upload/download failures for specific provider.
- **Impact**: Partial availability degradation.
- **Diagnosis Steps**:
  1. Check provider error HTTP response codes (e.g., 403 Forbidden, 429 Too Many Requests, 503 SlowDown).
- **Remediation Steps**:
  1. If 429 Too Many Requests, increase exponential backoff in retry configuration.
- **Escalation Path**: Page Cloud Storage Vendor Support.

---

## 7. StorageQuotaThresholdWarning (Medium)

- **Description**: Workspace storage quota usage > 85% for 15 minutes.
- **Symptoms**: Workspace is approaching max upload capacity.
- **Impact**: Uploads will be rejected once quota reaches 100%.
- **Remediation Steps**:
  1. Notify workspace administrator.
  2. Optionally increase workspace quota limit or trigger lifecycle expiration rules.
- **Escalation Path**: Tenant Account Administrator.

---

## 8. PolicyEngineFallbackRateHigh (Medium)

- **Description**: More than 2% of routing decisions are relying on default fallback routes (`cloudstorex_policy_fallback_routes_total`).
- **Symptoms**: Objects are stored in default provider rather than rule-specific provider.
- **Diagnosis Steps**:
  1. Inspect policy evaluation logs for condition evaluation errors.
- **Remediation Steps**:
  1. Correct syntax of recently updated routing rules.
- **Escalation Path**: Platform Policy Administrator.

---

## 9. MetadataCatalogGrowthWarning (Low)

- **Description**: Object catalog is growing faster than 1,000 objects per hour.
- **Remediation Steps**:
  1. Verify PostgreSQL table partitioning and index health.
- **Escalation Path**: DBA / SRE On-Call.
