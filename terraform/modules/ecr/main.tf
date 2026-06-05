# Server 用 Repository
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
