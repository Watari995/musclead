# ECS Cluster 本体(ただの入れ物)
resource "aws_ecs_cluster" "main" {
  name = "musclead-cluster"

  setting {
    name  = "containerInsights"
    value = "disabled"
  }

  tags = {
    Name = "musclead-cluster"
  }
}

# 使える計算リソースの宣言 + デフォルト
resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name       = aws_ecs_cluster.main.name
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    base              = 0
    weight            = 100
    capacity_provider = "FARGATE_SPOT"
  }
}

resource "aws_iam_role" "server_task_execution" {
  name = "musclead-server-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })

  tags = {
    Name = "musclead-server-task-execution-role"
  }
}

# AWS 標準セット: ECR pull + CloudWatch Logs 書き込みを許可
# (これは定番、 AWS が用意してる「ECS Task 起動の標準権限」)
resource "aws_iam_role_policy_attachment" "server_task_execution_basic" {
  role       = aws_iam_role.server_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# 独自追加: SSM Parameter Store の特定 ARN だけ読める権限
# (Task Definition の secrets 参照に必要、 最小権限で絞る)
resource "aws_iam_role_policy" "server_task_execution_ssm" {
  name = "musclead-server-ssm-read"
  role = aws_iam_role.server_task_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["ssm:GetParameters", "ssm:GetParameter"]
      Resource = var.ssm_parameter_arns # ←渡された ARN リストだけ許可
    }]
  })
}

# ─────────────────────────────────────────────────────────
# Task Role: コンテナ内の Go プロセスが AWS API を呼ぶための身分証
# (Task Execution Role とは別物、 § AWS の公式分離パターン)
# 用途: S3 PutObject / GetObject / DeleteObject(presigned URL 発行など)
# ─────────────────────────────────────────────────────────
resource "aws_iam_role" "server_task" {
  name = "musclead-server-task-role"

  # 「この Role は ECS Task が引き受ける」 と宣言
  # principal は ecs-tasks.amazonaws.com(Task Execution Role と同じ)
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })

  tags = {
    Name = "musclead-server-task-role"
  }
}

# S3 操作権限を Task Role に attach
# Resource を bucket arn 配下に絞ることで最小権限を維持
resource "aws_iam_role_policy" "server_task_storage" {
  name = "musclead-server-storage-access"
  role = aws_iam_role.server_task.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:PutObject",    # presigned PUT URL 発行(実際の upload はブラウザ→S3)
        "s3:GetObject",    # presigned GET URL 発行(画像表示用)
        "s3:DeleteObject", # 古い image を差し替え時に削除
      ]
      # bucket 配下のすべての object 対象、 bucket 自体の操作は不可
      Resource = "${var.storage_bucket_arn}/*"
    }]
  })
}

# CloudWatch Log Group: Server containerのログ集約先
resource "aws_cloudwatch_log_group" "server" {
  name              = "/musclead/ecs/server"
  retention_in_days = 7

  tags = {
    Name = "musclead-server-log"
  }
}


# Task Definition [containerを1つ起動する仕様書]をFargate用に作る
resource "aws_ecs_task_definition" "server" {
  family                   = "musclead-server"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"

  cpu    = "256"
  memory = "512"

  execution_role_arn = aws_iam_role.server_task_execution.arn
  # Task Role: コンテナ内 Go プロセスが AWS API(S3 等) を呼ぶための権限
  task_role_arn = aws_iam_role.server_task.arn

  runtime_platform {
    cpu_architecture        = "ARM64"
    operating_system_family = "LINUX"
  }

  # Container 定義: 1 Task に 1 つの Container
  # environment / secrets / logConfiguration は全て JSON の中に入れる
  container_definitions = jsonencode([{
    name      = "server"
    image     = var.server_image_url
    essential = true

    portMappings = [{
      containerPort = 8080
      protocol      = "tcp"
    }]

    # 平文 env(number は string に変換が必須)
    environment = [
      { name = "ADDR", value = ":8080" },
      { name = "DB_PORT", value = tostring(var.db_port) },
      { name = "DB_NAME", value = var.db_name },
      { name = "ALLOWED_ORIGIN", value = var.allowed_origin },
      { name = "AWS_REGION", value = var.aws_region },
      { name = "STORAGE_BUCKET", value = var.storage_bucket_name },
    ]

    # SSM 由来 secrets(IAM Role 経由で復号化)
    secrets = [
      { name = "JWT_SECRET", valueFrom = var.jwt_secret_arn },
      { name = "DB_USER", valueFrom = var.db_user_arn },
      { name = "DB_PASSWORD", valueFrom = var.db_password_arn },
      { name = "DB_HOST", valueFrom = var.db_host_arn },
    ]

    # CloudWatch Logs に awslogs ドライバで送信
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = aws_cloudwatch_log_group.server.name
        awslogs-region        = "ap-northeast-1"
        awslogs-stream-prefix = "server"
      }
    }
  }])

  tags = {
    Name = "musclead-server-task-definition"
  }
}

# ECS Service: Taskを起動するためのサービスを作る
resource "aws_ecs_service" "server" {
  name = "musclead-server-service"
  # どのClusterで動かすか
  cluster = aws_ecs_cluster.main.id
  # どのTask Definitionを使うか
  task_definition = aws_ecs_task_definition.server.arn
  # 常に1つのTaskを起動する
  desired_count = 1

  # 計算リソース: Fargate Spot 100%で起動する
  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 100
    base              = 0
  }

  # Taskの置き場所 + ネットワーク設定
  network_configuration {
    # どの Subnet にTaskをおくか
    subnets = var.subnet_ids

    # taskのセキュリティグループを指定
    security_groups = [var.server_sg_id]

    # Public IPを付与
    # ECR pull / SSM読み / Logの書き込みが外部API呼び出しなので
    assign_public_ip = true
  }

  # ALB Target Group に Task の IP を自動登録
  # Service が起動した Task の IP を Target Group に登録 → ALB が振り分け対象に
  load_balancer {
    target_group_arn = var.target_group_arn
    container_name   = "server" # Task Definition 内の container 名と一致
    container_port   = 8080
  }

  # IAM Roleのpolicy attach 完了を待ってからServiceを作る(race condition対策)
  depends_on = [aws_iam_role_policy_attachment.server_task_execution_basic, aws_iam_role_policy.server_task_execution_ssm]

  tags = {
    Name = "musclead-server-service"
  }
}
