variable "sqs_queue_arn" {
  description = "outbox SQS の ARN (event source mapping の source)"
  type        = string
}

variable "domain_name" {
  description = "SES ドメイン検証するドメイン (例: musclead.com)"
  type        = string
}

variable "hosted_zone_id" {
  description = "Route53 hosted zone ID (DKIM CNAME 作成用)"
  type        = string
}

variable "from_address" {
  description = "送信元アドレス (例: no-reply@musclead.com)"
  type        = string
}

variable "lambda_zip_path" {
  description = "ビルド済み Lambda zip のパス (bootstrap を zip 化したもの)"
  type        = string
}
