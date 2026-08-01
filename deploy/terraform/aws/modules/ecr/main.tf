resource "aws_ecr_repository" "backend" {
  name                 = "${var.environment}-cloudstorex-backend"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name        = "${var.environment}-cloudstorex-backend"
    Environment = var.environment
  }
}

resource "aws_ecr_repository" "frontend" {
  name                 = "${var.environment}-cloudstorex-frontend"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name        = "${var.environment}-cloudstorex-frontend"
    Environment = var.environment
  }
}
