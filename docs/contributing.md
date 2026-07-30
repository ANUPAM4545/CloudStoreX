# Contributing to CloudStoreX

Thank you for your interest in contributing to **CloudStoreX**! This document outlines our engineering workflow and review standards.

---

## 1. Branching Strategy

- **`main`**: The primary development branch. All pull requests target `main`.
- **Feature Branches**: Prefix branch names with the scope:
  - `feat/add-gcs-provider`
  - `fix/multipart-upload-timeout`
  - `docs/improve-architecture-diagram`

---

## 2. Commit Message Conventions

We adhere strictly to [Conventional Commits](https://www.conventionalcommits.org/):
```
<type>(<scope>): <short description>
```
**Types**:
- `feat`: A new feature or provider integration
- `fix`: A bug fix
- `docs`: Documentation-only changes
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding or updating tests
- `chore`: Maintenance on build scripts, CI, or dependencies

---

## 3. Pull Request Checklist

Before submitting a Pull Request, verify that:
1. All unit tests pass (`./scripts/test-all.sh`).
2. Code is cleanly formatted (`gofmt -s -w .` for Go, `npm run lint` for TypeScript).
3. No `TODO` or debug `console.log` statements are included.
4. New features include corresponding unit tests and documentation updates.
