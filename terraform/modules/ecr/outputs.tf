output "be_repository_url" {
  description = "BE Repository の URL(docker push / ECS Task Definition で参照)"
  value       = aws_ecr_repository.be.repository_url
}

output "be_repository_arn" {
  description = "BE Repository の ARN"
  value       = aws_ecr_repository.be.arn
}
