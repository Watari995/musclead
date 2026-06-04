output "endpoint" {
  description = "RDS の接続 endpoint(host)"
  value       = aws_db_instance.main.address
}

output "port" {
  description = "RDS の接続 port"
  value       = aws_db_instance.main.port
}

output "db_instance_identifier" {
  description = "RDS インスタンスの識別子"
  value       = aws_db_instance.main.identifier
}
