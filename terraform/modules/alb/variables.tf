variable "vpc_id" {
  description = "VPCのID"
  type        = string
}

variable "subnet_ids" {
  description = "ALBを配置するSubnet IDリスト(network moduleのpublic_subnet_ids)"
  type        = list(string)
}

variable "alb_sg_id" {
  description = "ALB用のSGのID(network moduleのalb_sg_id)"
  type        = string
}

variable "acm_certificate_arn" {
  description = "ACMで発行済の証明書ARN"
  type        = string
}
