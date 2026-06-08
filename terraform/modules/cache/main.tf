# Subnet Group : 2AZにまたがる必要 (RDSと同じ理由)
resource "aws_elasticache_subnet_group" "main" {
  name        = "musclead-cache-subnet-group"
  description = "musclead ElastiCache subnet group"
  subnet_ids  = var.subnet_ids

  tags = {
    Name = "musclead-cache-subnet-group"
  }
}

# Parameter Group : familyはengine_versionに合わせる
resource "aws_elasticache_parameter_group" "main" {
  name        = "musclead-cache-pg"
  family      = "redis7"
  description = "musclead Redis 7 parameter group"

  tags = {
    Name = "musclead-cache-pg"
  }
}

# Cluster: 単一ノードRedis
resource "aws_elasticache_replication_group" "main" {
  replication_group_id = "musclead-cache"
  description          = "musclead Redis replication group"
  engine               = "redis"
  engine_version       = "7.1"
  node_type            = "cache.t4g.micro"
  num_cache_clusters   = 1
  parameter_group_name = aws_elasticache_parameter_group.main.name
  port                 = 6379

  subnet_group_name  = aws_elasticache_subnet_group.main.name
  security_group_ids = [var.cache_sg_id] # network側でECSからのみ許可を定義する

  apply_immediately = true

  tags = {
    Name = "musclead-cache"
  }
}
