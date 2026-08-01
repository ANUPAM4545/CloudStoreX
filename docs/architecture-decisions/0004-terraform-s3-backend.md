# ADR 0004: Remote Terraform State Storage with S3 and DynamoDB Locking

## Status
Accepted (Epic 11)

## Context
Terraform state files (`.tfstate`) contain sensitive resource metadata and infrastructure mappings. Managing Terraform state locally or in git repositories exposes credentials and causes race conditions when multiple engineers or CI/CD pipelines attempt concurrent infrastructure updates.

## Decision
We standardize on remote AWS S3 backend storage with DynamoDB state locking for all Terraform deployments (`deploy/terraform/aws/`):
- **S3 Bucket**: Stores encrypted (`encrypt = true`), versioned state files with strict IAM access policies.
- **DynamoDB Table**: Provides state locking (`dynamodb_table = "cloudstorex-terraform-locks"`) to guarantee mutual exclusion during `terraform apply` operations.

## Consequences
### Positive
- **Concurrent Safety**: Prevents corrupt state or overwrites from concurrent runs.
- **Disaster Recovery**: S3 bucket versioning enables rollback to previous state snapshots if state corruption occurs.
- **Compliance**: State encryption at rest satisfies enterprise data protection policies.

### Negative / Mitigations
- **Bootstrap Requirement**: Requires an initial bootstrap S3 bucket and DynamoDB table before running `terraform init` for the first time (documented in `backend.tf.example`).
