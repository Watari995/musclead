variable "vpc_id" {
  description = "VPC ID(network module の output)"
  type        = string
}

variable "subnet_id" {
  description = "Bastion を置く Subnet ID(public subnet を 1 つ)"
  type        = string
}

variable "rds_sg_id" {
  description = "RDS の SG ID(ここに 3306 inbound を追加)"
  type        = string
}
