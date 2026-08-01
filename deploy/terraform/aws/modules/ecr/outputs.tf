output "backend_repository_url" {
  value       = aws_ecr_repository.backend.repository_url
  description = "URL of backend ECR repository"
}

output "frontend_repository_url" {
  value       = aws_ecr_repository.frontend.repository_url
  description = "URL of frontend ECR repository"
}
