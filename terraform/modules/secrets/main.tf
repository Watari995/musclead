# SSM Parameter Store(Systems Manager Parameter Store / 設定値保管庫)に
# 秘密値を SecureString として保管する。 ECS Task Definition から ARN 参照で読まれる。
#
# 値は terraform.tfvars で投入(Git 管理外)。 KMS は AWS default(aws/ssm)を使用、 無料。

resource "aws_ssm_parameter" "jwt_secret" {
  name        = "/musclead/${var.env}/jwt_secret"
  description = "JWT 署名鍵"
  type        = "SecureString"
  value       = var.jwt_secret

  tags = {
    Name = "musclead-jwt-secret"
  }
}

resource "aws_ssm_parameter" "db_user" {
  name        = "/musclead/${var.env}/db_user"
  description = "RDS MySQL ユーザー名"
  type        = "SecureString"
  value       = var.db_user

  tags = {
    Name = "musclead-db-user"
  }
}

resource "aws_ssm_parameter" "db_password" {
  name        = "/musclead/${var.env}/db_password"
  description = "RDS MySQL パスワード"
  type        = "SecureString"
  value       = var.db_password

  tags = {
    Name = "musclead-db-password"
  }
}

# db_host は RDS module 完成後に追加(RDS.endpoint output を value に渡す)
