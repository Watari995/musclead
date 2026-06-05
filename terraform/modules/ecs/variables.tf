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

# 平文 env(これはsecretではないので直接渡す)
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
