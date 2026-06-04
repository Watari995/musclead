variable "hosted_zone_id" {
  description = "Route 53 hosted zone ID(musclead.com の zone)"
  type        = string
}

variable "domain_name" {
  description = "ベースドメイン(例: musclead.com)"
  type        = string
}

variable "alb_dns_name" {
  description = "ALB の DNS 名(alb module の output)"
  type        = string
}

variable "alb_zone_id" {
  description = "ALB の Hosted Zone ID(Alias の評価先、 alb module の output)"
  type        = string
}
