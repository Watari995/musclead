output "cluster_name" {
  description = "ECS Clusterの名前"
  value       = aws_ecs_cluster.main.name
}

output "cluster_arn" {
  description = "ECS ClusterのARN"
  value       = aws_ecs_cluster.main.arn
}

output "be_task_execution_role_arn" {
  description = "BE Task Execution RoleのARN"
  value       = aws_iam_role.be_task_execution.arn
}

output "be_task_definition_arn" {
  description = "BE Task DefinitionのARN"
  value       = aws_ecs_task_definition.be.arn
}

output "be_service_name" {
  description = "BE ECS Service の名前(将来 ALB Target Group attach や CLI 操作で使う)"
  value       = aws_ecs_service.be.name
}
