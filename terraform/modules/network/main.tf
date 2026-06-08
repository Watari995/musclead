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

# ALB用 SG: 外からHTTPSを受け付ける
resource "aws_security_group" "alb" {
  name        = "musclead-alb-sg"
  description = "ALB: allow HTTPS from internet"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTPS from internet"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    description = "All outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = {
    Name = "musclead-alb-sg"
  }
}

#Server Fargate用 SG: ALB SGからのみ :8080を受け付ける
resource "aws_security_group" "server_fargate" {
  name        = "musclead-server-fargate-sg"
  description = "Server Fargate: allow :8080 from ALB SG only"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "App port from ALB SG"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }
  egress {
    description = "All outbound (ECR pull, RDS, etc.)"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
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

  ingress {
    description     = "MySQL from Server Fargate SG"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.server_fargate.id]
  }
  #egressはなし

  tags = {
    Name = "musclead-rds-sg"
  }
}

resource "aws_security_group" "cache" {
  name        = "musclead-cache-sg"
  description = "Cache: allow :6379 from BE Fargate SG only"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "Redis from BE Fargate SG"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.server_fargate.id]
  }
  egress {
    description = "All outbound (ECR pull, RDS, etc.)"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = {
    Name = "musclead-cache-sg"
  }
}
