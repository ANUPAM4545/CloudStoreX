# CloudStoreX Storage Domain Architecture

The **Storage Domain** is the heart of CloudStoreX. Every bucket creation, object upload, download, and deletion passes through this domain.

---

## 1. Core Domain Components

### `StorageService`
The `StorageService` interface (`internal/storage/service.go`) serves as the application-facing facade for all storage operations:
- Coordinates between input validators, the `StorageRouter`, and the response formatter.
- Encapsulates domain errors (`ErrBucketNotFound`, `ErrObjectNotFound`, `ErrProviderError`) and prevents backend provider exceptions from leaking to HTTP handlers.

### `StorageRouter`
The `StorageRouter` (`internal/storage/router.go`) is responsible for selecting the appropriate storage provider for a given request:
- Consults the `PolicyEvaluator` to inspect workspace policies, bucket tags, and tenant preferences.
- Retrieves the active `StorageProvider` instance from the `ProviderRegistry`.
- Directs the storage payload to the chosen provider without exposing routing complexity to the caller.

### `PolicyEvaluator`
The `PolicyEvaluator` (`internal/storage/policy.go`) evaluates declarative orchestration rules:
- Currently defaults to the `default` workspace policy (`MinIO` target provider).
- Designed to be extended with multi-cloud routing rules (e.g. lowest-cost routing, EU data residency, cold-storage replication).
