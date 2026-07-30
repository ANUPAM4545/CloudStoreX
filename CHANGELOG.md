# Changelog

All notable changes to **CloudStoreX** will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.1.0-alpha] - 2026-07-30

### Added

#### Epic 1: Platform Foundation
- **Monorepo Layout**: Clean top-level separation (`backend/`, `frontend/`, `docs/`, `docker/`, `scripts/`, `shared/`).
- **Go Backend Core**: Initialized Go 1.22+ project structure with Gin Web Framework, GORM, PostgreSQL 15, and Redis 7.
- **Enterprise Security**: Implemented JWT authentication and Argon2id/bcrypt password hashing.
- **Infrastructure Automation**: Created `docker-compose.yml` for zero-configuration local stack orchestration and `Makefile` automation.
- **Observability**: Added structured JSON logging (`internal/shared/logger`) and standardized API error responses.

#### Epic 2: Storage Domain & Provider Abstraction
- **Control Plane Architecture**: Implemented centralized `StorageService`, `StorageRouter`, and `PolicyEvaluator`.
- **Interface Segregation**: Created `StorageProvider`, `BucketProvider`, and `ObjectProvider` interfaces to decouple application logic from physical cloud backends.
- **Provider Registry**: Developed thread-safe runtime `Registry` using `sync.RWMutex` to support dynamic provider discovery and hot-swapping.
- **Policy Evaluator**: Created initial policy evaluation layer defaulting to standard MinIO workspace storage rules.

#### Epic 3: MinIO Provider Integration
- **MinIO Driver**: Implemented production-ready S3-compatible MinIO `StorageProvider` using `github.com/minio/minio-go/v7`.
- **Streaming I/O**: Designed zero-memory-buffer streaming `io.Reader`/`io.Writer` pipeline for high-throughput uploads and downloads.
- **Container Integration**: Integrated MinIO server and console into `docker-compose.yml` with healthchecks.

#### Epic 4: Storage REST API
- **REST Control Plane**: Exposed `/api/v1/storage` endpoints for managing buckets and objects.
- **Request Correlation**: Automatically injected `X-Request-ID` correlation header into all API responses.
- **Error Translation**: Implemented domain-to-HTTP status code mapping with structured JSON error schemas.
- **OpenAPI Docs**: Integrated code-generated OpenAPI documentation for storage endpoints.

#### Epic 5: Web Dashboard (Storage Management UI)
- **Modern Next.js 15 App**: Built enterprise web dashboard using Next.js 15 (App Router), TypeScript, Tailwind CSS, and shadcn/ui.
- **Feature Modules**: Organized codebase into clean self-contained feature modules (`auth`, `buckets`, `dashboard`, `objects`, `uploads`).
- **Hierarchical Object Explorer**: Built hierarchical folder navigation (`ObjectBreadcrumbs`, `ObjectTable`) simulating a standard filesystem over flat object keys.
- **Streaming Upload Manager**: Implemented floating bottom-right upload progress queue (`UploadProgress`) with individual `AbortController` cancellation and 100 MB upload limit safeguards.
