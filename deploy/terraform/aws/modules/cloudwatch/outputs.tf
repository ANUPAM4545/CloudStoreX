output "backend_log_group_name" {
  value       = aws_cloudwatch_log_group.backend.name
  description = "Name of backend log group"
}

output "frontend_log_group_name" {
  value       = aws_cloudwatch_log_group.frontend.name
  description = "Name of frontend log group"
}
