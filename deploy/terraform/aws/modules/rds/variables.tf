variable "environment" {
  type        = string
  description = "Deployment environment name"
}

variable "subnet_ids" {
  type        = list(string)
  description = "Private subnet IDs for DB subnet group"
}

variable "security_group_id" {
  type        = string
  description = "Security group ID for RDS instance"
}

variable "instance_class" {
  type        = string
  default     = "db.t3.medium"
  description = "RDS instance class"
}

variable "allocated_storage" {
  type        = number
  default     = 20
  description = "Allocated storage in GB"
}

variable "db_username" {
  type        = string
  default     = "cloudstorex"
  description = "Database admin username"
}

variable "db_password" {
  type        = string
  description = "Database admin password"
  sensitive   = true
}
