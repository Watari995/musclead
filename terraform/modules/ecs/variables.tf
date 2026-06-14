variable "ssm_parameter_arns" {
  description = "Taskが読む SSM ParameterのARNリスト(secrets moduleのoutputから渡す)"
  type        = list(string)
}

# Server Containerのimage URL
variable "server_image_url" {
  description = "Server Containerのimage URL (ex: 1234567890.dkr.ecr.ap-northeast-1.amazonaws.com/musclead-server:latest)"
  type        = string
}

# SSMからsecretsとして読む4つのARN
variable "jwt_secret_arn" {
  type      = string
  sensitive = true
}

variable "db_user_arn" {
  type      = string
  sensitive = true
}

variable "db_password_arn" {
  type      = string
  sensitive = true
}

variable "db_host_arn" {
  type      = string
  sensitive = true
}

variable "stripe_secret_key_arn" {
  type      = string
  sensitive = true
}

variable "stripe_webhook_signing_secret_arn" {
  type      = string
  sensitive = true
}

# 平文 env(これはsecretではないので直接渡す)

variable "stripe_pro_price_id" {
  description = "Stripe Pro プランの Price ID (price_)"
  type        = string
}

variable "stripe_success_url" {
  description = "Checkout 成功時の戻り URL"
  type        = string
}

variable "stripe_cancel_url" {
  description = "Checkout キャンセル時の戻り URL"
  type        = string
}

variable "stripe_portal_return_url" {
  description = "Customer Portal の戻り URL"
  type        = string
}
variable "db_name" {
  type = string
}

variable "db_port" {
  type = number
}

variable "allowed_origin" {
  description = "CORS 許可オリジン(FE の URL)"
  type        = string
}

variable "subnet_ids" {
  description = "Taskを配置するSubnet IDリスト(network moduleのpublic_subnet_ids)"
  type        = list(string)
}

variable "server_sg_id" {
  description = "Server Fargate用SGのID(network moduleのserver_fargate_sg_id)"
  type        = string
}

variable "target_group_arn" {
  description = "ALB Target Group の ARN(alb moduleの output)"
  type        = string
}

# ── Storage(S3 image bucket)─────────────

variable "storage_bucket_name" {
  description = "S3 bucket 名(storage module の output、 BE の env として渡す)"
  type        = string
}

variable "storage_bucket_arn" {
  description = "S3 bucket ARN(Task Role の S3 policy で参照)"
  type        = string
}

variable "aws_region" {
  description = "AWS region(BE が AWS SDK 初期化で使用)"
  type        = string
}


# ── Cache ─────────────
variable "cache_endpoint" {
  description = "Cacheのエンドポイント"
  type        = string
}
