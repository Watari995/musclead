variable "domain_name" {
  description = "証明書を発行するメインドメイン(例: musclead.com)"
  type        = string
}

variable "hosted_zone_id" {
  description = "Route 53 hosted zone ID(検証用 CNAME を投入する先)"
  type        = string
}
