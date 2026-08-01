output "app_role_arn" {
  value       = aws_iam_role.app.arn
  description = "ARN of the app IAM role"
}

output "app_role_name" {
  value       = aws_iam_role.app.name
  description = "Name of the app IAM role"
}
