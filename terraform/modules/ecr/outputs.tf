output "server_repository_url" {
  description = "Server Repository の URL(docker push / ECS Task Definition で参照)"
  value       = aws_ecr_repository.server.repository_url
}

output "server_repository_arn" {
  description = "Server Repository の ARN"
  value       = aws_ecr_repository.server.arn
}
