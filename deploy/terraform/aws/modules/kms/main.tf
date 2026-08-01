# KMS module placeholder (Refinement 9)
# Future enhancement: Dedicated KMS customer-managed keys for S3 SSE-KMS, RDS, and ElastiCache encryption at rest.

variable "environment" {
  type        = string
  description = "Deployment environment name"
  default     = "dev"
}
