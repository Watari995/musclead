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

module "network" {
  source = "./modules/network"
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

module "ecs" {
  source       = "./modules/ecs"
  be_image_url = "${module.ecr.be_repository_url}:latest"
  ssm_parameter_arns = [
    module.secrets.jwt_secret_arn,
    module.secrets.db_user_arn,
    module.secrets.db_password_arn,
    module.secrets.db_host_arn,
  ]
  jwt_secret_arn  = module.secrets.jwt_secret_arn
  db_user_arn     = module.secrets.db_user_arn
  db_password_arn = module.secrets.db_password_arn
  db_host_arn     = module.secrets.db_host_arn
  db_name         = var.db_name
  db_port         = var.db_port
}


# module "alb" {
#   source = "./modules/alb"
#   ...
# }

# module "dns" {
#   source      = "./modules/dns"
#   domain_name = var.domain_name
#   alb_dns     = module.alb.alb_dns_name
#   alb_zone_id = module.alb.alb_zone_id
# }
