output "alb_dns_name" {
  description = "ALBのDNS名"
  value       = aws_lb.main.dns_name
}

output "be_target_group_arn" {
  description = "BE Target GroupのARN"
  value       = aws_lb_target_group.be.arn
}

output "alb_zone_id" {
  description = "ALBのHosted Zone ID(Route 53 Alias レコードで参照)"
  value       = aws_lb.main.zone_id
}
