# outbox consumer (Lambda) 一式。 SQS → Lambda → DynamoDB(冪等) → SES (ADR 0020)。
# 全部 VPC 外 / pay-per-use = 追加費用ゼロ。

# === DynamoDB: 通知の重複排除テーブル (TTL で自動削除) ===
resource "aws_dynamodb_table" "dedup" {
  name         = "musclead-outbox-dedup"
  billing_mode = "PAY_PER_REQUEST" # オンデマンド (低トラフィック = ほぼ $0)
  hash_key     = "event_id"

  attribute {
    name = "event_id"
    type = "S"
  }

  ttl {
    attribute_name = "expires_at"
    enabled        = true
  }
}

# === Lambda 実行ロール ===
resource "aws_iam_role" "consumer" {
  name = "musclead-outbox-consumer"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name = "musclead-outbox-consumer"
  }
}

# 基本実行 (CloudWatch Logs への書き込み)
resource "aws_iam_role_policy_attachment" "consumer_basic" {
  role       = aws_iam_role.consumer.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# SQS 受信 (event source mapping が Lambda にメッセージを渡すのに必要)
resource "aws_iam_role_policy_attachment" "consumer_sqs" {
  role       = aws_iam_role.consumer.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
}

# アプリ権限: DynamoDB (PutItem/DeleteItem) + SES (SendEmail)
resource "aws_iam_role_policy" "consumer_app" {
  name = "musclead-outbox-consumer-app"
  role = aws_iam_role.consumer.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:PutItem",    # 冪等用の条件付き書き込み
          "dynamodb:DeleteItem", # SES 失敗時の印削除
        ]
        Resource = aws_dynamodb_table.dedup.arn # この dedup テーブルだけに限定
      },
      {
        Effect   = "Allow"
        Action   = ["ses:SendEmail"]
        Resource = "*" # SES の送信はリソース指定が粗いので "*" が一般的
      },
    ]
  })
}

# === Lambda 関数本体 ===
resource "aws_lambda_function" "consumer" {
  function_name = "musclead-outbox-consumer"
  role          = aws_iam_role.consumer.arn
  runtime       = "provided.al2023" # Go は custom runtime (bootstrap バイナリ)
  handler       = "bootstrap"
  architectures = ["arm64"]

  filename         = var.lambda_zip_path                   # 動かすコード (bootstrap を zip 化したもの)
  source_code_hash = filebase64sha256(var.lambda_zip_path) # zip の中身が変わった時だけ再デプロイ

  environment {
    variables = {
      DEDUP_TABLE_NAME = aws_dynamodb_table.dedup.name
      SES_FROM_ADDRESS = var.from_address
    }
  }
}

# === SQS → Lambda トリガー (= ②の event source mapping) ===
resource "aws_lambda_event_source_mapping" "outbox" {
  event_source_arn                   = var.sqs_queue_arn          # どの SQS が
  function_name                      = aws_lambda_function.consumer.arn # どの Lambda を呼ぶか
  batch_size                         = 10                         # 1 回の呼び出しで最大 10 件まとめて渡す
  maximum_batching_window_in_seconds = 5                          # 最大 5 秒は溜めてからまとめて呼ぶ (呼び出し回数を減らす)
  # enabled = true がデフォルト (作った瞬間から有効)
}
