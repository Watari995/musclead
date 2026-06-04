variable "env" {
  description = "環境名(SSM パス階層に使う、 例: prod)"
  type        = string
}

variable "jwt_secret" {
  description = "JWT 署名鍵"
  type        = string
  sensitive   = true
}

variable "db_user" {
  description = "RDS MySQL ユーザー名"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "RDS MySQL パスワード"
  type        = string
  sensitive   = true
}
