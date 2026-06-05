# 踏み台 EC2: TablePlus → SSM トンネル → ここ → RDS の中継役
#
# 設計:
#  - t4g.nano (ARM, $0.0042/hr)
#  - public subnet 配置、 public IP あり (SSM Agent が AWS API に届くため)
#  - SG inbound = なし (SSM は完全 outbound)
#  - IAM: AmazonSSMManagedInstanceCore のみ
#  - 普段は stopped、 mus-prod alias で start/stop

# 最新のAmazon Linux 2023 ARM AMIを動的取得
data "aws_ami" "al2023_arm" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-arm64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# IAM Role: EC2 が AWS サービスを呼ぶための「身分証」
resource "aws_iam_role" "bastion" {
  name = "musclead-bastion-role"

  # 「EC2 サービスがこの Role を装着してよい」 という宣言
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

# IAM Role: AmazonSSMManagedInstanceCore のみ
resource "aws_iam_role_policy_attachment" "bastion_ssm" {
  role       = aws_iam_role.bastion.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "bastion" {
  name = "musclead-bastion-profile"
  role = aws_iam_role.bastion.name
}

resource "aws_security_group" "bastion" {
  name        = "musclead-bastion-sg"
  description = "Bastion EC2 for SSM tunneling to RDS"

  #どのVPCのSGか (network moduleから渡される)
  vpc_id = var.vpc_id

  # outbound ルール(送信側)
  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"

    # 宛先 IP: 0.0.0.0/0 (全てのIPに対して全てのプロトコルを許可)
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "musclead-bastion-sg"
  }
}

# RDSのSGにBastionからの3306 inboundを許可する
resource "aws_security_group_rule" "rds_from_bastion" {
  type      = "ingress"
  from_port = 3306
  to_port   = 3306
  protocol  = "tcp"
  # どのSGに対してこのRULEを入れるか (対象はRDSのSG)
  security_group_id = var.rds_sg_id
  # どのSGからの通信を指定する
  source_security_group_id = aws_security_group.bastion.id

  description = "Allow Bastion to access RDS"
}


# Bastion EC2 本体
resource "aws_instance" "bastion" {
  ami = data.aws_ami.al2023_arm.id
  # インスタンスサイズ
  instance_type = "t4g.nano"
  subnet_id     = var.subnet_id

  # SGリスト
  vpc_security_group_ids = [aws_security_group.bastion.id]

  # IAM Profile = このEC2がSSM Agentを起動できるようになる
  iam_instance_profile = aws_iam_instance_profile.bastion.name

  # Public IPを付与
  associate_public_ip_address = true

  # ルートディスク
  root_block_device {
    volume_size = 8
    volume_type = "gp3"
    encrypted   = true
  }

  tags = {
    Name = "musclead-bastion"
  }
}
