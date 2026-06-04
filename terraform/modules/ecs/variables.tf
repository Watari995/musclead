variable "ssm_parameter_arns" {
  description = "Taskが読む SSM ParameterのARNリスト(secrets moduleのoutputから渡す)"
  type        = list(string)
}

# BE Containerのimage URL
variable "be_image_url" {
  description = "BE Containerのimage URL (ex: 1234567890.dkr.ecr.ap-northeast-1.amazonaws.com/musclead-be:latest)"
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
