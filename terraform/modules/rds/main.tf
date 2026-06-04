# RDSを置くSubnet群 (2AZにまたがる必要、AWSの規約)
resource "aws_db_subnet_group" "main" {
  name        = "musclead-rds-subnet-group"
  description = "musclead RDS subnet group"
  subnet_ids  = var.subnet_ids

  tags = {
    Name = "musclead-rds-subnet-group"
  }
}

resource "aws_db_parameter_group" "main" {
  name        = "musclead-rds-pg"
  family      = "mysql8.0"
  description = "musclead RDS MySQL 8.0 parameter group"

  tags = {
    Name = "musclead-rds-pg"
  }
}

resource "aws_db_instance" "main" {
  identifier     = "musclead-rds"
  engine         = "mysql"
  engine_version = "8.0"
  instance_class = "db.t4g.micro"

  // ディスクの容量
  allocated_storage     = 20
  max_allocated_storage = 20
  storage_type          = "gp3"
  storage_encrypted     = true

  db_name  = var.db_name
  username = var.db_user
  password = var.db_password
  port     = var.db_port

  db_subnet_group_name   = aws_db_subnet_group.main.name
  parameter_group_name   = aws_db_parameter_group.main.name
  vpc_security_group_ids = [var.rds_sg_id]
  publicly_accessible    = false

  multi_az                = false
  backup_retention_period = 1
  skip_final_snapshot     = true
  deletion_protection     = false

  tags = {
    Name = "musclead-rds"
  }
}
