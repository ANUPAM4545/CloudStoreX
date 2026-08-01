output "endpoint" {
  value       = aws_elasticache_cluster.this.cache_nodes[0].address
  description = "Hostname endpoint of the ElastiCache Redis cluster"
}

output "port" {
  value       = aws_elasticache_cluster.this.cache_nodes[0].port
  description = "Port of the ElastiCache Redis cluster"
}
