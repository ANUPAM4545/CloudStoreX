variable "environment" {
  type        = string
  description = "Deployment environment name"
}

variable "bucket_suffix" {
  type        = string
  description = "Unique suffix for S3 bucket name"
  default     = "001"
}

variable "force_destroy" {
  type        = bool
  description = "Allow destroying non-empty S3 buckets in non-prod environments"
  default     = false
}
