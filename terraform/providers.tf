provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project   = "musclead"
      ManagedBy = "Terraform"
      Env       = var.env
    }
  }
}
