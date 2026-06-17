# === SES: ドメイン検証 + DKIM 署名 ===
# ドメイン検証にすると no-reply@<domain> 等から送れ、 DKIM 署名で到達率が上がる (ADR 0020)。

resource "aws_sesv2_email_identity" "domain" {
  email_identity = var.domain_name
  # Easy DKIM はデフォルトで有効。 dkim_signing_attributes.tokens に検証用トークン(3個)が入る
}

# DKIM 検証用の CNAME を Route53 に3本作る (これで SES がドメイン所有を確認できる)。
# Easy DKIM は必ず token を3個返す。 for_each はキーが値依存で apply 前に確定できないため、
# 個数が静的に決まる count = 3 を使う (トークンの値だけ apply 後に解決される)。
resource "aws_route53_record" "dkim" {
  count = 3

  zone_id = var.hosted_zone_id
  name    = "${aws_sesv2_email_identity.domain.dkim_signing_attributes[0].tokens[count.index]}._domainkey.${var.domain_name}"
  type    = "CNAME"
  ttl     = 300
  records = ["${aws_sesv2_email_identity.domain.dkim_signing_attributes[0].tokens[count.index]}.dkim.amazonses.com"]
}

# 補足: 新規 SES アカウントは「サンドボックス」 (検証済み宛先のみ送信可)。
#       一般ユーザーに送るには AWS に production access を申請する (Terraform 外の手作業)。
