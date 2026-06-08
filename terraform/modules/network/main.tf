resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "musclead-vpc"
  }
}

resource "aws_subnet" "public_1a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "ap-northeast-1a"
  map_public_ip_on_launch = true

  tags = {
    Name = "musclead-public-1a"
  }
}

resource "aws_subnet" "public_1c" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = "ap-northeast-1c"
  map_public_ip_on_launch = true

  tags = {
    Name = "musclead-public-1c"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "musclead-igw"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "musclead-public-rt"
  }
}

resource "aws_route" "public_internet" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.main.id
}

resource "aws_route_table_association" "public_1a" {
  subnet_id      = aws_subnet.public_1a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_1c" {
  subnet_id      = aws_subnet.public_1c.id
  route_table_id = aws_route_table.public.id
}

# ─────────────────────────────────────────────────────────
# Security Groups (本体)
# Rule は全て standalone (aws_security_group_rule) で別途定義する。
# inline ingress/egress と standalone rule は混在禁止 (Terraform 公式)。
# musclead は cross-module で rule 追加するため standalone に統一。
# ─────────────────────────────────────────────────────────

resource "aws_security_group" "alb" {
  name        = "musclead-alb-sg"
  description = "ALB: allow HTTPS from internet"
  vpc_id      = aws_vpc.main.id
  tags = {
    Name = "musclead-alb-sg"
  }
}

resource "aws_security_group" "server_fargate" {
  name        = "musclead-server-fargate-sg"
  description = "Server Fargate: allow :8080 from ALB SG only"
  vpc_id      = aws_vpc.main.id
  tags = {
    Name = "musclead-server-fargate-sg"
  }
}

resource "aws_security_group" "rds" {
  name = "musclead-rds-sg"
  # SG description は作成時しか設定できない(変更すると SG 再作成 → ENI cleanup で失敗)
  # BE 命名で作成済のため、 機能変更なら別 SG 新設、 ここは旧文字列維持
  description = "RDS: allow :3306 from BE Fargate SG only"
  vpc_id      = aws_vpc.main.id
  tags = {
    Name = "musclead-rds-sg"
  }
}

resource "aws_security_group" "cache" {
  name        = "musclead-cache-sg"
  description = "Cache: allow :6379 from BE Fargate SG only"
  vpc_id      = aws_vpc.main.id
  tags = {
    Name = "musclead-cache-sg"
  }
}

# ─────────────────────────────────────────────────────────
# Security Group Rules (standalone)
# ─────────────────────────────────────────────────────────

# ALB ingress
resource "aws_security_group_rule" "alb_https_in" {
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  security_group_id = aws_security_group.alb.id
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "HTTPS from internet"
}

# ALB egress
resource "aws_security_group_rule" "alb_egress_all" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  security_group_id = aws_security_group.alb.id
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "All outbound"
}

# Server Fargate ingress (ALB SG からの 8080)
resource "aws_security_group_rule" "server_fargate_app_in" {
  type                     = "ingress"
  from_port                = 8080
  to_port                  = 8080
  protocol                 = "tcp"
  security_group_id        = aws_security_group.server_fargate.id
  source_security_group_id = aws_security_group.alb.id
  description              = "App port from ALB SG"
}

# Server Fargate egress (ECR pull, RDS, ElastiCache 等の outbound 用)
resource "aws_security_group_rule" "server_fargate_egress_all" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  security_group_id = aws_security_group.server_fargate.id
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "All outbound (ECR pull, RDS, ElastiCache, etc.)"
}

# RDS ingress (Server Fargate SG からの 3306)
resource "aws_security_group_rule" "rds_mysql_in" {
  type                     = "ingress"
  from_port                = 3306
  to_port                  = 3306
  protocol                 = "tcp"
  security_group_id        = aws_security_group.rds.id
  source_security_group_id = aws_security_group.server_fargate.id
  description              = "MySQL from Server Fargate SG"
}

# Cache ingress (Server Fargate SG からの 6379)
resource "aws_security_group_rule" "cache_redis_in" {
  type                     = "ingress"
  from_port                = 6379
  to_port                  = 6379
  protocol                 = "tcp"
  security_group_id        = aws_security_group.cache.id
  source_security_group_id = aws_security_group.server_fargate.id
  description              = "Redis from Server Fargate SG"
}

# Cache egress (Redis は受動的だが念のため open)
resource "aws_security_group_rule" "cache_egress_all" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  security_group_id = aws_security_group.cache.id
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "All outbound (rarely needed for Redis; allow for safety)"
}
