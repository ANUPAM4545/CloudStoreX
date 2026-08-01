variable "aws_region" {
  type        = string
  description = "AWS region"
  default     = "us-east-1"
}

variable "environment" {
  type        = string
  description = "Deployment environment (dev, staging, prod)"
  default     = "dev"
}

variable "vpc_cidr" {
  type        = string
  default     = "10.0.0.0/16"
  description = "CIDR block for VPC"
}

variable "availability_zones" {
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
  description = "Availability zones for subnet allocation"
}

variable "public_subnet_cidrs" {
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
  description = "CIDR blocks for public subnets"
}

variable "private_subnet_cidrs" {
  type        = list(string)
  default     = ["10.0.10.0/24", "10.0.11.0/24"]
  description = "CIDR blocks for private subnets"
}

variable "db_instance_class" {
  type        = string
  default     = "db.t3.medium"
  description = "RDS instance class"
}

variable "db_allocated_storage" {
  type        = number
  default     = 20
  description = "RDS allocated storage in GB"
}

variable "db_password" {
  type        = string
  description = "Database administrator password"
  sensitive   = true
}

variable "redis_node_type" {
  type        = string
  default     = "cache.t3.micro"
  description = "ElastiCache Redis node type"
}

variable "s3_bucket_suffix" {
  type        = string
  default     = "001"
  description = "Unique suffix for data S3 bucket"
}

variable "log_retention_days" {
  type        = number
  default     = 30
  description = "CloudWatch log retention in days"
}
