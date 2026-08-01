variable "environment" {
  type        = string
  description = "Deployment environment name"
}

variable "s3_bucket_arn" {
  type        = string
  description = "ARN of the S3 bucket for data storage"
}
