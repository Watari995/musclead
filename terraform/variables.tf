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

# サブドメイン:
# - app.<domain_name>: FE
# - api.<domain_name>: BE
