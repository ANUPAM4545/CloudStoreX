# CloudStoreX Terraform AWS Infrastructure Reference

The Terraform configuration under `deploy/terraform/aws/` provisions production AWS cloud infrastructure for CloudStoreX (**Refinement 9**).

---

## 1. Modules

- `modules/vpc/`: VPC, public/private subnets, IGW, routing tables.
- `modules/security-groups/`: App security group (HTTP/Frontend) and DB security group (PostgreSQL 5432 / Redis 6379 ingress).
- `modules/rds/`: Amazon RDS PostgreSQL instance in private subnet group.
- `modules/redis/`: Amazon ElastiCache Redis cluster in private subnet group.
- `modules/s3/`: Versioned, AES256-encrypted data bucket with public access block enabled.
- `modules/iam/`: Application IAM Role and least-privilege policy for S3 bucket access.
- `modules/ecr/`: Immutable ECR repositories for backend and frontend container images.
- `modules/cloudwatch/`: CloudWatch log groups (`/aws/cloudstorex/{env}/backend` and `frontend`) with configurable retention.
- `modules/kms/`, `route53/`, `acm/`, `alb/`: Architectural placeholder modules for enterprise encryption and traffic routing extensions.

---

## 2. Remote State & Locking (Refinement 3)
In production, state is stored in an encrypted AWS S3 backend with DynamoDB locking:
```hcl
# backend.tf
terraform {
  backend "s3" {
    bucket         = "cloudstorex-terraform-state-prod"
    key            = "cloudstorex/production/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "cloudstorex-terraform-locks"
  }
}
```
