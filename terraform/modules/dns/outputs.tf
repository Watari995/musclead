output "api_fqdn" {
  description = "BE API の FQDN(api.musclead.com)"
  value       = aws_route53_record.api.fqdn
}
