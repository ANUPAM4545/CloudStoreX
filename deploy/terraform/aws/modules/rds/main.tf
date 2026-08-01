resource "aws_db_subnet_group" "this" {
  name       = "${var.environment}-cloudstorex-rds-subnet-group"
  subnet_ids = var.subnet_ids

  tags = {
    Name        = "${var.environment}-cloudstorex-rds-subnet-group"
    Environment = var.environment
  }
}

resource "aws_db_instance" "this" {
  identifier             = "${var.environment}-cloudstorex-postgres"
  engine                 = "postgres"
  engine_version         = "16.1"
  instance_class         = var.instance_class
  allocated_storage      = var.allocated_storage
  storage_type           = "gp3"
  db_name                = "cloudstorex"
  username               = var.db_username
  password               = var.db_password
  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [var.security_group_id]
  skip_final_snapshot    = true
  publicly_accessible    = false

  tags = {
    Name        = "${var.environment}-cloudstorex-postgres"
    Environment = var.environment
  }
}
