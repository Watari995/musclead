
variable "subnet_ids" {
  description = "ElastiCacheを配置するSubnet IDs (network moduleのpublic_subnet_ids)"
  type        = list(string)
}

variable "cache_sg_id" {
  description = "ElastiCache用のSGのID(network module出力)"
  type        = string
}
