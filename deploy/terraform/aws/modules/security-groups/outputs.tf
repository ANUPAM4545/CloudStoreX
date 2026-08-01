output "app_security_group_id" {
  value       = aws_security_group.app.id
  description = "ID of the app security group"
}

output "db_security_group_id" {
  value       = aws_security_group.db.id
  description = "ID of the database security group"
}
