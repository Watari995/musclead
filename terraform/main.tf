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

# module "network" {
#   source = "./modules/network"
#   env    = var.env
# }

# module "rds" {
#   source           = "./modules/rds"
#   env              = var.env
#   vpc_id           = module.network.vpc_id
#   subnet_ids       = module.network.public_subnet_ids
#   db_sg_id         = module.network.db_sg_id
# }

# module "ecr" {
#   source = "./modules/ecr"
#   env    = var.env
# }

# module "ecs" {
#   source = "./modules/ecs"
#   ...
# }

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
