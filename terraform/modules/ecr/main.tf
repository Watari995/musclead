# Server 用 Repository
# FE は Vercel デプロイに変更したため、 ECR は BE 1 個だけ(ADR 0007 改訂 / 0008 参照)
resource "aws_ecr_repository" "server" {
  name                 = "musclead-server"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "musclead-server"
  }
}
