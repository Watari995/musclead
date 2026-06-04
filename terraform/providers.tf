provider "aws" {
  region  = var.aws_region
  profile = "musclead-admin"

  default_tags {
    tags = {
      Project   = "musclead"
      ManagedBy = "Terraform"
      Env       = var.env
    }
  }
}
