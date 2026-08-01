# ADR 0001: Use Distroless Images for Production Go Backend Runtime

## Status
Accepted (Epic 11)

## Context
CloudStoreX is a high-security, cloud-native storage control plane that processes authenticated user credentials and object metadata. Standard base images (such as `golang:1.23` or standard Debian/Ubuntu runtimes) contain package managers (`apt`, `dpkg`), shells (`bash`, `sh`), and OS utilities that increase container image size and present an expanded attack surface for arbitrary code execution or privilege escalation if a vulnerability is exploited.

## Decision
We have standardized on `gcr.io/distroless/static-debian12:nonroot` as the production runtime image for the backend service (`cloudstorex-backend`).
- The multi-stage build compiles a static binary (`CGO_ENABLED=0`) in a `golang:1.23-alpine` stage.
- The compiled binary is copied into the minimal distroless runtime container running as UID `65532` (`nonroot`).

## Consequences
### Positive
- **Minimal Attack Surface**: Zero shells or package managers; prevents interactive shell execution or container escape payloads.
- **Minimal Image Size**: The runtime image is under ~30MB, significantly reducing ECR storage costs, image pull times, and Kubernetes node startup latency.
- **Non-Root Execution**: Runs strictly as an unprivileged user (`UID 65532`), complying with Pod Security Standards (Restricted).

### Negative / Mitigations
- **No Interactive Debugging**: `kubectl exec` into the running container with a shell is not possible. Debugging requires ephemeral debug containers (`kubectl debug`) or structured CloudWatch logs and metrics.
