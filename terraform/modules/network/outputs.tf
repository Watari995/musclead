output "vpc_id" {
  description = "VPCのID"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "Public SubnetのIDリスト(複数AZ)"
  value       = [aws_subnet.public_1a.id, aws_subnet.public_1c.id]
}

output "alb_sg_id" {
  description = "ALB用のSGのID"
  value       = aws_security_group.alb.id
}

output "be_fargate_sg_id" {
  description = "BE Fargate用SGのID"
  value       = aws_security_group.be_fargate.id
}

output "rds_sg_id" {
  description = "RDS用SGのID"
  value       = aws_security_group.rds.id
}
