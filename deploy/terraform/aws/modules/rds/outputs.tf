output "endpoint" {
  value       = aws_db_instance.this.endpoint
  description = "Endpoint of the created RDS instance"
}

output "db_name" {
  value       = aws_db_instance.this.db_name
  description = "Database name"
}
