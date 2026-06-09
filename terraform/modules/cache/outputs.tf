output "endpoint" {
  description = "ElastiCacheのエンドポイント"
  value       = aws_elasticache_replication_group.main.primary_endpoint_address
}

output "port" {
  description = "ElastiCacheのポート"
  value       = aws_elasticache_replication_group.main.port
}
