variable "subnet_ids" {
  description = "RDSを配置するSubnet IDリスト(network moduleのpublic_subnet_ids)"
  type        = list(string)
}

variable "rds_sg_id" {
  description = "RDS 用 SG の ID(network module の出力)"
  type        = string
}

variable "db_name" {
  description = "MySQL データベース名"
  type        = string
}

variable "db_user" {
  description = "MySQL ユーザー名"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "MySQL パスワード"
  type        = string
  sensitive   = true
}

variable "db_port" {
  description = "MySQL 接続ポート"
  type        = number
}
