# ADR 0003: Conditional Embedded vs. External Dependencies in Helm Chart

## Status
Accepted (Epic 11)

## Context
CloudStoreX requires PostgreSQL, Redis, and S3-compatible object storage. For local Kubernetes development (e.g., Kind, Minikube) and automated testing environments, developers require a self-contained chart that deploys all dependencies without manual cloud resource provisioning. However, in enterprise production environments, embedded StatefulSets lack cross-AZ replication, automated point-in-time backups, and managed SLA guarantees.

## Decision
We implemented conditional dependency flags in `helm/cloudstorex`:
- Embedded `postgresql`, `redis`, and `minio` sub-charts/templates can be toggled via `postgresql.enabled=true/false`, `redis.enabled=true/false`, and `minio.enabled=true/false`.
- Dynamic configuration helper templates (`_helpers.tpl`) resolve database, Redis, and storage endpoints:
  - When embedded is enabled, services point to internal Kubernetes cluster DNS names (`cloudstorex-postgresql`, `cloudstorex-redis`, `cloudstorex-minio`).
  - When disabled, services inject external host endpoints (`externalDatabase.host`, `externalRedis.host`, `externalStorage.endpoint`) and external secrets via environment variables.

## Consequences
### Positive
- **Single Chart for All Environments**: The exact same chart (`helm/cloudstorex`) is deployed across dev, staging, and production without maintaining separate manifests.
- **Zero Risk of Accidental Embedded Prod Deployment**: Clear parameter overrides in `values-prod.yaml` enforce external AWS RDS, ElastiCache, and S3 usage.

### Negative / Mitigations
- **Values Complexity**: Operators must ensure proper external credentials are provided when turning off embedded dependencies. Documented in `docs/kubernetes.md`.
