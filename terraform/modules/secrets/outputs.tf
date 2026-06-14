# ECS Task Definition の secrets フィールドが参照する ARN を export。

output "jwt_secret_arn" {
  description = "JWT 署名鍵 Parameter の ARN"
  value       = aws_ssm_parameter.jwt_secret.arn
}

output "db_user_arn" {
  description = "DB ユーザー名 Parameter の ARN"
  value       = aws_ssm_parameter.db_user.arn
}

output "db_password_arn" {
  description = "DB パスワード Parameter の ARN"
  value       = aws_ssm_parameter.db_password.arn
}

output "db_host_arn" {
  description = "DB host Parameter の ARN"
  value       = aws_ssm_parameter.db_host.arn
}

output "stripe_secret_key_arn" {
  description = "Stripe secret key Parameter の ARN"
  value       = aws_ssm_parameter.stripe_secret_key.arn
}

output "stripe_webhook_signing_secret_arn" {
  description = "Stripe Webhook 署名 secret Parameter の ARN"
  value       = aws_ssm_parameter.stripe_webhook_signing_secret.arn
}
