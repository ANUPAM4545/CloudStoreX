resource "aws_cloudwatch_log_group" "backend" {
  name              = "/aws/cloudstorex/${var.environment}/backend"
  retention_in_days = var.retention_in_days

  tags = {
    Name        = "${var.environment}-cloudstorex-backend-logs"
    Environment = var.environment
  }
}

resource "aws_cloudwatch_log_group" "frontend" {
  name              = "/aws/cloudstorex/${var.environment}/frontend"
  retention_in_days = var.retention_in_days

  tags = {
    Name        = "${var.environment}-cloudstorex-frontend-logs"
    Environment = var.environment
  }
}
