# CloudStoreX Production Security Hardening Guide & Audit Report

This document details the security architecture, controls, and hardening measures implemented in **Epic 11** for production enterprise deployments of CloudStoreX.

---

## 1. Container Security & Non-Root Execution
- **Distroless Runtime**: The CloudStoreX backend executes in `gcr.io/distroless/static-debian12:nonroot`, which removes interactive shells, operating system tools, and package managers.
- **UID / GID Restriction**:
  - Backend pods run as UID `65532` (`nonroot`).
  - Frontend pods run as UID `1001` (`nextjs`).
- **Kubernetes SecurityContext**:
  - `runAsNonRoot: true`: Prevents the Kubernetes kubelet from starting any container if the image runs as UID 0.
  - `allowPrivilegeEscalation: false`: Denies any process inside the container from gaining additional privileges via `setuid`/`setgid`.
  - `readOnlyRootFilesystem: true` (backend): Protects against runtime file modification attacks. Scratch disk writes are restricted to an ephemeral `/tmp` `emptyDir` mount.
  - `seccompProfile.type: RuntimeDefault`: Restricts kernel syscalls to standard safe defaults.

---

## 2. Automated Vulnerability Management (CI/CD)
- **Software Bill of Materials (SBOM)**: Every container image build automatically generates an SPDX-JSON SBOM using `anchore/sbom-action` (powered by Syft).
- **Vulnerability Scanning**: `aquasecurity/trivy-action` scans images for critical and high-severity OS and library vulnerabilities before pushing to the container registry.
- **Keyless Image Signing**: `sigstore/cosign` cryptographically signs all container images pushed to GHCR using OIDC identity tokens, enabling admission controllers (e.g., Kyverno or Sigstore Policy Controller) to verify image authenticity before deployment.
- **Automated Dependency Maintenance**: Dependabot (`.github/dependabot.yml`) weekly audits Go modules, NPM packages, GitHub Actions workflows, and Docker base images for upstream CVE patches.

---

## 3. Kubernetes Network Segmentation
- **NetworkPolicy Enforcement**: `cloudstorex-network-policy` applies default-deny ingress rules within the `cloudstorex` namespace.
- **Allowed Communication Paths**:
  - Ingress Controller -> Frontend Service (`3000`) & Backend Service (`8080`).
  - Frontend Pods -> Backend Service (`8080`).
  - Backend Pods -> PostgreSQL (`5432`), Redis (`6379`), MinIO/External Storage (`9000`).
  - All other unauthorized pod-to-pod or external ingress connections are dropped.

---

## 4. Secret Management & TLS Termination
- **Secret Separation**: Sensitive credentials (`DB_PASSWORD`, `JWT_SECRET`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`) are decoupled from ConfigMaps and injected purely through Kubernetes Secrets or ExternalSecrets operators (integrating with AWS Secrets Manager / HashiCorp Vault).
- **TLS Termination**: Ingress rules (`ingress.yaml`) support TLS termination via cert-manager or AWS ACM certificate ARNs.
- **JWT Key Hardening**: Production deployments enforce 256-bit+ cryptographic HMAC/RSA keys for workspace token signing and validation.
