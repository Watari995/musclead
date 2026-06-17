# outbox relay 用の SQS。 ECS の relay worker が publish し、 Lambda consumer が処理する。
# 無料枠 (100万 req/月) 内 = 追加費用ゼロ (ADR 0020)。

# DLQ: consumer が maxReceiveCount 回処理に失敗したメッセージを退避する先 (毒メッセージ対策)。
resource "aws_sqs_queue" "outbox_dlq" {
  name                      = "musclead-outbox-dlq"
  message_retention_seconds = 1209600 # 14日

  tags = {
    Name = "musclead-outbox-dlq"
  }
}

# メインキュー: outbox relay が publish する先。
resource "aws_sqs_queue" "outbox" {
  name                       = "musclead-outbox"
  visibility_timeout_seconds = 30     # consumer の処理猶予 (Lambda の実行時間に合わせる)
  message_retention_seconds  = 345600 # 4日

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.outbox_dlq.arn
    maxReceiveCount     = 5 # 5回失敗で DLQ へ
  })

  tags = {
    Name = "musclead-outbox"
  }
}
