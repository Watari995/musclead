# tfstate を S3 に保存、 DynamoDB でロックする本番品質構成。
# S3 バケットと DynamoDB テーブルは手動で先に作成しておく必要がある
# (循環依存を避けるため Terraform 管理外)。
#
# bucket / dynamodb_table の名前は実際に作成したものに合わせて変更すること。

terraform {
  required_version = "~> 1.9.8"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "musclead-tfstate-204340689570"
    key            = "terraform.tfstate"
    region         = "ap-northeast-1"
    dynamodb_table = "musclead-tfstate-lock"
    encrypt        = true
  }
}
