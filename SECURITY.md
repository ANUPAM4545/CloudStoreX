# Security Policy – CloudStoreX

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.11.x  | :white_check_mark: |
| < 1.10  | :x:                |

## Reporting a Vulnerability

We take the security of CloudStoreX and its multi-cloud storage control plane seriously.

If you discover a security vulnerability within CloudStoreX, please follow these steps:
1. **Do not disclose publicly**: Do not create a public GitHub Issue for potential vulnerabilities.
2. **Contact Security Team**: Email `security@cloudstorex.company.com` with:
   - A description of the vulnerability and its potential impact.
   - Detailed steps to reproduce the issue (proof-of-concept scripts or HTTP traces are appreciated).
   - Any suggested mitigations.
3. **Response Timeline**:
   - You will receive an acknowledgment within **24 hours**.
   - A preliminary assessment and remediation plan will be provided within **5 business days**.

## Security Hardening Practices
CloudStoreX enforces strict security practices across its deployment lifecycle:
- **Non-Root Container Execution**: All container images run as unprivileged non-root users (`65532` or `1001`).
- **Minimal Runtime Images**: Backend runs on Google Distroless static images without shells or package managers.
- **Automated Scanning**: CI/CD pipelines run SBOM generation (Syft) and vulnerability scanning (Trivy) on every build.
- **Image Signing**: All published images are signed via Sigstore Cosign.
- **Network Isolation**: Default Kubernetes NetworkPolicies restrict inter-pod communication and deny unauthorized ingress.
