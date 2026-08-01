# CloudStoreX Enterprise Helm Chart Reference

The `cloudstorex` Helm chart (`helm/cloudstorex`) deploys the CloudStoreX backend, frontend, ingress, autoscaling, and optional embedded data services on Kubernetes.

---

## 1. Values Configuration Parameters

| Parameter | Default | Description |
| :--- | :--- | :--- |
| `global.environment` | `"development"` | Target environment name (`development`, `staging`, `production`) |
| `backend.replicaCount` | `1` | Number of backend container replicas (overridden by HPA in prod) |
| `backend.image.repository` | `"cloudstorex-backend"` | Backend Docker image repository |
| `backend.image.tag` | `"latest"` | Backend Docker image tag |
| `postgresql.enabled` | `true` | Deploy embedded PostgreSQL StatefulSet (set `false` for prod) |
| `redis.enabled` | `true` | Deploy embedded Redis StatefulSet (set `false` for prod) |
| `minio.enabled` | `true` | Deploy embedded MinIO storage (set `false` for prod) |
| `externalDatabase.host` | `""` | Hostname of external AWS RDS PostgreSQL instance |
| `externalRedis.host` | `""` | Hostname of external AWS ElastiCache Redis cluster |
| `externalStorage.bucket` | `""` | Target AWS S3 bucket name |

---

## 2. Deploying with External Dependencies
Example command for enterprise production environments:
```bash
helm upgrade --install cloudstorex-prod ./helm/cloudstorex \
  --namespace cloudstorex-prod --create-namespace \
  --set postgresql.enabled=false \
  --set redis.enabled=false \
  --set minio.enabled=false \
  --set externalDatabase.host="prod-postgres.cloudstorex.internal" \
  --set externalRedis.host="prod-redis.cloudstorex.internal" \
  --set externalStorage.bucket="company-cloudstorex-prod-bucket" \
  --set autoscaling.enabled=true \
  --set autoscaling.minReplicas=2 \
  --set autoscaling.maxReplicas=10
```
