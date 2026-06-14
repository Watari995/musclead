variable "aws_region" {
  description = "デプロイ先 AWS リージョン"
  type        = string
  default     = "ap-northeast-1"
}

variable "env" {
  description = "環境名 (prod / staging 等)"
  type        = string
  default     = "prod"
}

variable "domain_name" {
  description = "登録済みドメイン名 (Route 53 で取得)"
  type        = string
  default     = "musclead.com"
}

# ── Secrets / 環境ごとに変わる値(全て terraform.tfvars で投入) ─────────────

variable "jwt_secret" {
  description = "JWT 署名鍵(SSM Parameter Store の SecureString に保管)"
  type        = string
  sensitive   = true
}

variable "db_user" {
  description = "RDS MySQL マスターユーザー名"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "RDS MySQL マスターパスワード"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "RDS MySQL のスキーマ名"
  type        = string
}

variable "db_port" {
  description = "RDS MySQL の接続ポート"
  type        = number
}

# ── Route 53 / ACM ─────────────

variable "hosted_zone_id" {
  description = "Route 53 hosted zone ID(musclead.com の zone)"
  type        = string
}

# ── CORS ─────────────

variable "allowed_origin" {
  description = "Server が CORS 許可するオリジン(FE の URL)"
  type        = string
  default     = "https://app.musclead.com"
}

# ── Cache ─────────────
variable "enable_cache" {
  description = "value_cacheを有効にするかどうか"
  type        = bool
  default     = false
}

# ── Stripe ─────────────
variable "stripe_secret_key" {
  description = "Stripe API secret key (sk_test / sk_live)。 terraform.tfvars で投入"
  type        = string
  sensitive   = true
}

variable "stripe_webhook_signing_secret" {
  description = "Stripe Webhook 署名検証 secret (whsec_)。 terraform.tfvars で投入"
  type        = string
  sensitive   = true
}

variable "stripe_pro_price_id" {
  description = "Stripe Pro プランの Price ID"
  type        = string
  default     = "price_1TgMut866skNoey5280aSIKe"
}

variable "stripe_success_url" {
  description = "Checkout 成功時の戻り URL"
  type        = string
  default     = "https://app.musclead.com/settings?purchase=success"
}

variable "stripe_cancel_url" {
  description = "Checkout キャンセル時の戻り URL"
  type        = string
  default     = "https://app.musclead.com/settings?purchase=cancel"
}

variable "stripe_portal_return_url" {
  description = "Customer Portal の戻り URL"
  type        = string
  default     = "https://app.musclead.com/settings"
}
