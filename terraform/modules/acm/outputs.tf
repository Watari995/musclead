output "certificate_arn" {
  description = "発行済 ACM 証明書の ARN(ALB Listener の certificate_arn で参照)"
  # 注意: aws_acm_certificate.main.arn ではなく validation の方を参照する
  # → これによって「検証完了している証明書」 だけが出力される
  value = aws_acm_certificate_validation.main.certificate_arn
}
