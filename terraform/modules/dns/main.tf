# api.musclead.com → ALB に向ける Alias レコード
#
# なぜ CNAME ではなく Alias?
#   サブドメインなら CNAME でも動くが、 Alias は AWS リソース専用の上位互換:
#     - 無料(CNAME クエリは課金、 Alias は無料)
#     - 高速(AWS 内部で直接解決)
#     - root ドメインでも使える(CNAME は root に置けない)
#   ALB は AWS リソースなので Alias 一択。
resource "aws_route53_record" "api" {
  zone_id = var.hosted_zone_id
  name    = "api.${var.domain_name}" # api.musclead.com
  type    = "A"                      # Alias でも DNS 上は A レコード扱い

  alias {
    name                   = var.alb_dns_name # 例: musclead-alb-xxxx.ap-northeast-1.elb.amazonaws.com
    zone_id                = var.alb_zone_id  # ALB の Hosted Zone ID(ALB module の output)
    evaluate_target_health = true             # ALB が unhealthy なら DNS 解決失敗(無駄リクエスト減らす)
  }
}

# app.musclead.com → Vercel 用 A レコード
#
# Vercel は anycast IP (76.76.21.21) で全プロジェクトを受け、 SNI ヘッダで振り分け。
# CNAME (cname.vercel-dns.com) も動くが、 Vercel 公式は A レコード推奨。
resource "aws_route53_record" "app" {
  zone_id = var.hosted_zone_id
  name    = "app.${var.domain_name}" # app.musclead.com
  type    = "A"
  ttl     = 60
  records = ["76.76.21.21"] # Vercel anycast IP
}
