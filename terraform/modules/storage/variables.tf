variable "account_id" {
  description = "AWS Account ID (bucket 名の衝突回避用)"
  type        = string
}

variable "allowed_origins" {
  description = "CORS で許可するブラウザ origin リスト(FE の URL)"
  type        = list(string)
  default     = ["https://app.musclead.com"]
}
