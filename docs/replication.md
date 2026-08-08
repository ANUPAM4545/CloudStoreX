# Replication

CloudStoreX provides asynchronous and synchronous cross-region object replication to ensure high availability and durability.

## Architecture

Replication workers subscribe to `ObjectUploaded` events via the Reliability Event Bus and enqueue replication jobs. The system tracks `BytesReplicated`, lag, and object checksums for integrity.

## Metrics

- `cloudstorex_reliability_replication_bytes_total`
- `cloudstorex_reliability_replication_errors_total`
- `cloudstorex_reliability_replication_lag_seconds`
