output "alb_dns_name" {
  description = "ALBのDNS名"
  value       = aws_lb.main.dns_name
}

output "server_target_group_arn" {
  description = "Server Target GroupのARN"
  value       = aws_lb_target_group.server.arn
}

output "alb_zone_id" {
  description = "ALBのHosted Zone ID(Route 53 Alias レコードで参照)"
  value       = aws_lb.main.zone_id
}
