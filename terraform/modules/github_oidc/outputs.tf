output "role_arn" {
  description = "GitHub Actions が assume する Role の ARN(workflow の aws-actions/configure-aws-credentials に渡す)"
  value       = aws_iam_role.github_actions.arn
}
