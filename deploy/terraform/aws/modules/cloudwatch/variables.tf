variable "environment" {
  type        = string
  description = "Deployment environment name"
}

variable "retention_in_days" {
  type        = number
  default     = 30
  description = "CloudWatch logs retention in days"
}
