# CloudStoreX Disaster Recovery, Backup & Restore Procedures

This document outlines data preservation, snapshotting, and recovery procedures for CloudStoreX metadata and object payloads.

---

## 1. Recovery Point Objective (RPO) & Recovery Time Objective (RTO)
- **Metadata Database (PostgreSQL)**:
  - **RPO**: `< 5 minutes` via automated RDS point-in-time recovery (PITR) and WAL archiving.
  - **RTO**: `< 15 minutes` via automated RDS Multi-AZ failover or snapshot restoration.
- **Object Payloads (S3 / Storage Providers)**:
  - **RPO**: `Zero data loss` for committed objects (S3 11-nines durability + bucket versioning).
  - **RTO**: Instantaneous availability via provider multi-region replication.

---

## 2. PostgreSQL Backup & Restore

### Automated RDS PITR
In AWS production environments, RDS automated backups are enabled with a 30-day retention window.
- **Restore Point-in-Time**:
  ```bash
  aws rds restore-db-instance-to-point-in-time \
    --source-db-instance-identifier prod-cloudstorex-postgres \
    --target-db-instance-identifier prod-cloudstorex-postgres-restored \
    --restore-time 2026-08-01T06:00:00Z
  ```

### Manual Backup (Embedded Kubernetes DB)
To export PostgreSQL database dumps from the embedded Kubernetes pod:
```bash
# Backup
kubectl exec -it cloudstorex-postgresql-0 -n cloudstorex -- \
  pg_dump -U cloudstorex cloudstorex > cloudstorex_backup_$(date +%F).sql

# Restore
kubectl exec -i cloudstorex-postgresql-0 -n cloudstorex -- \
  psql -U cloudstorex cloudstorex < cloudstorex_backup_2026-08-01.sql
```
