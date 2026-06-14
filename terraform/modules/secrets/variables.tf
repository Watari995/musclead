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

variable "db_host" {
  description = "RDS endpoint(rds module の output から受け取る)"
  type        = string
}

variable "stripe_secret_key" {
  description = "Stripe API secret key (sk_test / sk_live)"
  type        = string
  sensitive   = true
}

variable "stripe_webhook_signing_secret" {
  description = "Stripe Webhook 署名検証 secret (whsec_)"
  type        = string
  sensitive   = true
}
