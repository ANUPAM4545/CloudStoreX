# CloudStoreX Docker Orchestration

This directory houses container orchestration assets, environment configurations, and custom deployment overrides for **CloudStoreX**.

## Architecture Overview

CloudStoreX's local development stack (`docker-compose.yml` in repository root) launches the following services:

| Service | Container Name | Port | Description |
|---|---|---|---|
| **PostgreSQL 15** | `cloudstorex_postgres` | `5432` | Relational metadata, user accounts, and workspace registry |
| **Redis 7** | `cloudstorex_redis` | `6379` | Distributed session caching and high-speed policy rate limiting |
| **MinIO** | `cloudstorex_minio` | `9000` (API), `9001` (Console) | Object storage provider target for local testing |
| **Backend (Go)** | `cloudstorex_backend` | `8080` | Storage Router, Provider Registry, and REST API Control Plane |
| **Frontend (Next.js)** | `cloudstorex_frontend` | `3000` | Enterprise Storage Management Web Dashboard |

## Environment Configuration

Use `docker/.env.example` as a template for overriding default environment variables in production or staging environments.
