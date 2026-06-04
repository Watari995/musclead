# ACM 証明書をリクエスト + 自動 DNS 検証 + 検証完了まで wait の 3 ステップを一括 IaC 化
#
# 流れ:
#   1. aws_acm_certificate: 「musclead.com + *.musclead.com の証明書ください」
#      → ACM が「この CNAME を投入してください」 と検証情報を返す
#   2. aws_route53_record (for_each): その CNAME を Route 53 に投入
#      → ACM が定期的にチェック、 ドメイン所有を確認
#   3. aws_acm_certificate_validation: 検証完了を Terraform が待つだけのリソース
#      → ALB が「Issued」 状態の証明書しか使えないため、 ここの完了が必須

resource "aws_acm_certificate" "main" {
  # メインのドメイン
  domain_name = var.domain_name

  # SAN(Subject Alternative Names): ワイルドカードで app/api 等のサブドメイン全部カバー
  subject_alternative_names = ["*.${var.domain_name}"]

  # 検証方式: DNS(Route 53 連携で完全自動。 Email 方式より楽 + 高速)
  validation_method = "DNS"

  # 安全対策: 古い証明書を消す前に新しいのを作る
  # (ALB 等が使ってる証明書を先に消すと一瞬 HTTPS 落ちるため)
  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "musclead-cert"
  }
}

# DNS 検証用レコードを Route 53 に投入
# ACM が「これを置け」 と指示する CNAME を、 dynamic に複数作る
# (domain_name + SAN それぞれに 1 レコード必要)
resource "aws_route53_record" "validation" {
  # for_each: ACM が返してくる検証情報(複数)を 1 つずつ展開
  for_each = {
    for dvo in aws_acm_certificate.main.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  # 既にレコードがあれば上書き(再 apply 時の安全策)
  allow_overwrite = true

  name    = each.value.name
  records = [each.value.record]
  ttl     = 60
  type    = each.value.type
  zone_id = var.hosted_zone_id
}

# 検証完了を待つ (実リソースは作らない、 「待ち」 専用の仮想リソース)
# Terraform 的にこのリソースの完了が「証明書が Issued になった」 を保証する
resource "aws_acm_certificate_validation" "main" {
  certificate_arn         = aws_acm_certificate.main.arn
  validation_record_fqdns = [for record in aws_route53_record.validation : record.fqdn]
}
