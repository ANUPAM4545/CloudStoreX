# CloudStoreX Local Development Guide

This guide covers setting up your local environment for developing and debugging **CloudStoreX**.

---

## 1. Prerequisites

- **Go**: `1.22+`
- **Node.js**: `v20+ LTS` and `npm`
- **Docker & Docker Compose**: `v2.20+`

---

## 2. Makefile Command Reference

The root `Makefile` automates common engineering workflows:

| Command | Description |
|---|---|
| `make up` | Launches the complete stack in Docker (Postgres, Redis, MinIO, Backend, Frontend) |
| `make down` | Stops and removes all Docker Compose containers |
| `make db-up` | Launches only local infrastructure (Postgres, Redis, MinIO) for host development |
| `make db-down` | Stops local infrastructure containers |
| `make build` | Builds Go server binary (`bin/server`) and Next.js frontend production bundle |
| `make test` | Runs Go backend unit and repository tests (`go test -v -cover ./...`) |
| `make lint` | Runs `golangci-lint` on backend and `npm run lint` on frontend |
| `make clean` | Removes compiled binaries and build artifacts |

---

## 3. Running Services Locally

1. Start Postgres, Redis, and MinIO:
   ```bash
   make db-up
   ```
2. Run the Go Backend API Server:
   ```bash
   make run
   ```
   The backend will listen on `http://localhost:8080`.
3. Run the Next.js Web Dashboard:
   ```bash
   cd frontend && npm run dev
   ```
   The dashboard will listen on `http://localhost:3000`.
