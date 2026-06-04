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
