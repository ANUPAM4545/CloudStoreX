output "vpc_id" {
  value       = module.vpc.vpc_id
  description = "ID of the created VPC"
}

output "rds_endpoint" {
  value       = module.rds.endpoint
  description = "RDS PostgreSQL endpoint"
}

output "redis_endpoint" {
  value       = module.redis.endpoint
  description = "ElastiCache Redis hostname endpoint"
}

output "s3_bucket_name" {
  value       = module.s3.bucket_name
  description = "Data S3 bucket name"
}

output "app_iam_role_arn" {
  value       = module.iam.app_role_arn
  description = "IAM Role ARN for application pods/instances"
}

output "backend_ecr_url" {
  value       = module.ecr.backend_repository_url
  description = "Backend ECR repository URL"
}

output "frontend_ecr_url" {
  value       = module.ecr.frontend_repository_url
  description = "Frontend ECR repository URL"
}
