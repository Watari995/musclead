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

resource "aws_iam_role" "be_task_execution" {
  name = "musclead-be-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })

  tags = {
    Name = "musclead-be-task-execution-role"
  }
}

# AWS 標準セット: ECR pull + CloudWatch Logs 書き込みを許可
# (これは定番、 AWS が用意してる「ECS Task 起動の標準権限」)
resource "aws_iam_role_policy_attachment" "be_task_execution_basic" {
  role       = aws_iam_role.be_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# 独自追加: SSM Parameter Store の特定 ARN だけ読める権限
# (Task Definition の secrets 参照に必要、 最小権限で絞る)
resource "aws_iam_role_policy" "be_task_execution_ssm" {
  name = "musclead-be-ssm-read"
  role = aws_iam_role.be_task_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["ssm:GetParameters", "ssm:GetParameter"]
      Resource = var.ssm_parameter_arns # ←渡された ARN リストだけ許可
    }]
  })
}

# CloudWatch Log Group: BE containerのログ集約先
resource "aws_cloudwatch_log_group" "be" {
  name              = "/musclead/ecs/be"
  retention_in_days = 7

  tags = {
    Name = "musclead-be-log"
  }
}


# Task Definition [containerを1つ起動する仕様書]をFargate用に作る
resource "aws_ecs_task_definition" "be" {
  family                   = "musclead-be"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"

  cpu    = "256"
  memory = "512"

  execution_role_arn = aws_iam_role.be_task_execution.arn

  # Container 定義: 1 Task に 1 つの Container
  # environment / secrets / logConfiguration は全て JSON の中に入れる
  container_definitions = jsonencode([{
    name      = "be"
    image     = var.be_image_url
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
        awslogs-group         = aws_cloudwatch_log_group.be.name
        awslogs-region        = "ap-northeast-1"
        awslogs-stream-prefix = "be"
      }
    }
  }])

  tags = {
    Name = "musclead-be-task-definition"
  }
}
