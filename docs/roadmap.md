# CloudStoreX Roadmap

This roadmap tracks the development lifecycle of **CloudStoreX** from Alpha toward production multi-cloud general availability.

---

## 1. Completed Milestones (`v0.1.0-alpha`)

- [x] **Epic 1: Platform Foundation**: Monorepo layout, Go + Gin + GORM + PostgreSQL + Redis, JWT auth, structured JSON logging, and Makefile automation.
- [x] **Epic 2: Storage Domain**: Centralized `StorageService`, `StorageRouter`, and `PolicyEvaluator`.
- [x] **Epic 3: MinIO Provider**: Pluggable S3-compatible MinIO provider driver with streaming I/O.
- [x] **Epic 4: Storage REST API**: `/api/v1/storage` endpoints with request ID correlation and OpenAPI docs.
- [x] **Epic 5: Web Dashboard**: Next.js 15 enterprise web dashboard with hierarchical folder navigation and upload queue.

---

## 2. In Progress (`v0.2.0-alpha`)

- [ ] **Metadata Engine**: Dedicated PostgreSQL/Redis object catalog supporting custom user metadata, object tags, and lightning-fast indexing.
- [ ] **Multi-Cloud Providers**: Native `StorageProvider` implementations for AWS S3, Google Cloud Storage (GCS), and Azure Blob Storage.

---

## 3. Upcoming Milestones (`v0.3.0-beta`)

- [ ] **Declarative Policy Engine**: YAML-based orchestration policies (cost optimization, region affinity, automated replication).
- [ ] **Storage Intelligence**: Automated cold-storage tiering, object lifecycle policies, and deduplication analytics.
