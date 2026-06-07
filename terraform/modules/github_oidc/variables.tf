variable "github_repo" {
  description = "GitHub リポジトリ <org>/<repo> 形式(例: Watari995/musclead)"
  type        = string
}

variable "allowed_branch" {
  description = "deploy を許可するブランチ名(他ブランチからは Role 引き受け不可)"
  type        = string
  default     = "main"
}

variable "ecr_repository_arn" {
  description = "ECR repository の ARN(push 先を絞るため)"
  type        = string
}

variable "task_execution_role_arn" {
  description = "ECS Task Execution Role の ARN(PassRole 対象)"
  type        = string
}

variable "task_role_arn" {
  description = "ECS Task Role の ARN(PassRole 対象、 アプリが AWS API を呼ぶための role)"
  type        = string
}
