variable "environment" {
  type        = string
  description = "Deployment environment name"
}

variable "subnet_ids" {
  type        = list(string)
  description = "Private subnet IDs for Redis subnet group"
}

variable "security_group_id" {
  type        = string
  description = "Security group ID for ElastiCache cluster"
}

variable "node_type" {
  type        = string
  default     = "cache.t3.micro"
  description = "ElastiCache node type"
}
