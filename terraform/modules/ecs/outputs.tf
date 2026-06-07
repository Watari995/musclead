output "cluster_name" {
  description = "ECS Clusterの名前"
  value       = aws_ecs_cluster.main.name
}

output "cluster_arn" {
  description = "ECS ClusterのARN"
  value       = aws_ecs_cluster.main.arn
}

output "server_task_execution_role_arn" {
  description = "Server Task Execution RoleのARN"
  value       = aws_iam_role.server_task_execution.arn
}

output "server_task_role_arn" {
  description = "Server Task Role の ARN(アプリ実行時の AWS API 権限)"
  value       = aws_iam_role.server_task.arn
}

output "server_task_definition_arn" {
  description = "Server Task DefinitionのARN"
  value       = aws_ecs_task_definition.server.arn
}

output "server_service_name" {
  description = "Server ECS Service の名前(将来 ALB Target Group attach や CLI 操作で使う)"
  value       = aws_ecs_service.server.name
}
