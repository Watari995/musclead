# 各 module の root レベル配線。 module の中身は user が順次実装する想定。
#
# 依存関係:
#   network → rds / ecs (subnet_id, sg を共有)
#   ecr     → ecs (image url を渡す)
#   ecs     → alb (target group attachment)
#   alb     → dns (Alias record target)
#   secrets → ecs (Task Definition で参照)
#
# 初期は network のみ apply して動作確認、 順次他を有効化していく。

# AWS Account ID を動的取得(S3 bucket 名の衝突回避用)
data "aws_caller_identity" "current" {}

module "network" {
  source = "./modules/network"
}

module "storage" {
  source          = "./modules/storage"
  account_id      = data.aws_caller_identity.current.account_id
  allowed_origins = ["https://app.${var.domain_name}"]
}

module "rds" {
  source      = "./modules/rds"
  subnet_ids  = module.network.public_subnet_ids
  rds_sg_id   = module.network.rds_sg_id
  db_name     = var.db_name
  db_user     = var.db_user
  db_password = var.db_password
  db_port     = var.db_port
}

module "ecr" {
  source = "./modules/ecr"
}

module "secrets" {
  source      = "./modules/secrets"
  env         = var.env
  jwt_secret  = var.jwt_secret
  db_user     = var.db_user
  db_password = var.db_password
  db_host     = module.rds.endpoint
}

module "acm" {
  source         = "./modules/acm"
  domain_name    = var.domain_name
  hosted_zone_id = var.hosted_zone_id
}

module "alb" {
  source              = "./modules/alb"
  vpc_id              = module.network.vpc_id
  subnet_ids          = module.network.public_subnet_ids
  alb_sg_id           = module.network.alb_sg_id
  acm_certificate_arn = module.acm.certificate_arn
}

module "ecs" {
  source           = "./modules/ecs"
  server_image_url = "${module.ecr.server_repository_url}:latest"
  ssm_parameter_arns = [
    module.secrets.jwt_secret_arn,
    module.secrets.db_user_arn,
    module.secrets.db_password_arn,
    module.secrets.db_host_arn,
  ]
  jwt_secret_arn      = module.secrets.jwt_secret_arn
  db_user_arn         = module.secrets.db_user_arn
  db_password_arn     = module.secrets.db_password_arn
  db_host_arn         = module.secrets.db_host_arn
  db_name             = var.db_name
  db_port             = var.db_port
  subnet_ids          = module.network.public_subnet_ids
  server_sg_id        = module.network.server_fargate_sg_id
  target_group_arn    = module.alb.server_target_group_arn
  allowed_origin      = var.allowed_origin
  storage_bucket_name = module.storage.bucket_name
  storage_bucket_arn  = module.storage.bucket_arn
  aws_region          = var.aws_region
}


module "dns" {
  source         = "./modules/dns"
  hosted_zone_id = var.hosted_zone_id
  domain_name    = var.domain_name
  alb_dns_name   = module.alb.alb_dns_name
  alb_zone_id    = module.alb.alb_zone_id
}

module "bastion" {
  source    = "./modules/bastion"
  vpc_id    = module.network.vpc_id
  subnet_id = module.network.public_subnet_ids[0]
  rds_sg_id = module.network.rds_sg_id
}

module "github_oidc" {
  source                  = "./modules/github_oidc"
  github_repo             = "Watari995/musclead"
  allowed_branch          = "main"
  ecr_repository_arn      = module.ecr.server_repository_arn
  task_execution_role_arn = module.ecs.server_task_execution_role_arn
  task_role_arn           = module.ecs.server_task_role_arn
}
