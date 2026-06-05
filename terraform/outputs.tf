output "github_actions_role_arn" {
  description = "GitHub Actions が assume する Role の ARN(workflow に貼り付ける)"
  value       = module.github_oidc.role_arn
}

output "ecr_repository_url" {
  description = "Server ECR Repository URL(workflow の docker push 先)"
  value       = module.ecr.server_repository_url
}

output "ecs_cluster_name" {
  description = "ECS Cluster 名(workflow の deploy 対象)"
  value       = module.ecs.cluster_name
}

output "ecs_service_name" {
  description = "ECS Service 名(workflow の deploy 対象)"
  value       = module.ecs.server_service_name
}
