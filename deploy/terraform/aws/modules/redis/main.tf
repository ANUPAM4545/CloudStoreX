resource "aws_elasticache_subnet_group" "this" {
  name       = "${var.environment}-cloudstorex-redis-subnet-group"
  subnet_ids = var.subnet_ids
}

resource "aws_elasticache_cluster" "this" {
  cluster_id           = "${var.environment}-cloudstorex-redis"
  engine               = "redis"
  node_type            = var.node_type
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.this.name
  security_group_ids   = [var.security_group_id]

  tags = {
    Name        = "${var.environment}-cloudstorex-redis"
    Environment = var.environment
  }
}
