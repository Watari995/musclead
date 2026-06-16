output "queue_url" {
  description = "outbox メインキューの URL (ECS の OUTBOX_QUEUE_URL に渡す)"
  value       = aws_sqs_queue.outbox.url
}

output "queue_arn" {
  description = "outbox メインキューの ARN (IAM / Lambda トリガーで使う)"
  value       = aws_sqs_queue.outbox.arn
}

output "dlq_arn" {
  description = "DLQ の ARN"
  value       = aws_sqs_queue.outbox_dlq.arn
}
