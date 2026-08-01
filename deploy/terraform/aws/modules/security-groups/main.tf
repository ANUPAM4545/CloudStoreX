resource "aws_security_group" "app" {
  name        = "${var.environment}-cloudstorex-app-sg"
  description = "Security group for CloudStoreX backend and frontend compute pods"
  vpc_id      = var.vpc_id

  ingress {
    description = "HTTP traffic from internal network"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  ingress {
    description = "Frontend traffic from internal network"
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.environment}-cloudstorex-app-sg"
    Environment = var.environment
  }
}

resource "aws_security_group" "db" {
  name        = "${var.environment}-cloudstorex-db-sg"
  description = "Security group for RDS PostgreSQL and ElastiCache Redis"
  vpc_id      = var.vpc_id

  ingress {
    description     = "PostgreSQL from app security group"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.app.id]
  }

  ingress {
    description     = "Redis from app security group"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.app.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.environment}-cloudstorex-db-sg"
    Environment = var.environment
  }
}
