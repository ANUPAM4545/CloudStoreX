# CloudStoreX Coding Standards

This document establishes the idiomatic coding rules for contributors to **CloudStoreX**.

---

## 1. Go Backend Standards (`backend/`)

- **Standard Project Layout**: Follow idiomatic Go structure (`cmd/`, `internal/...`). Never place application logic in `pkg/`.
- **Interface Segregation**: Keep interfaces small and consumer-focused (`BucketProvider`, `ObjectProvider`).
- **Error Handling**: Always return wrapped errors or use structured domain errors from `internal/shared/response`. Never panic in HTTP handlers or provider drivers.
- **Context Propagation**: Always pass `context.Context` as the first parameter to any database, cache, or storage provider operation.

---

## 2. TypeScript & Next.js Standards (`frontend/`)

- **Feature Modules**: Organize business logic under `src/features/<feature-name>/`. Keep `src/components/ui/` reserved for generic presentation primitives.
- **DTO to View Model Mapping**: Never consume raw backend REST DTOs inside React components. Transform API DTOs into View Models using dedicated mapper functions.
- **Design Tokens**: All components must consume centralized design tokens from `globals.css` (`--radius: 0.75rem`, color variables). Avoid ad-hoc color hex codes or random Tailwind spacing.
- **Strict Typing**: Avoid `any` types. Ensure all API request/response DTOs are strongly typed in `src/lib/api/types.ts`.
