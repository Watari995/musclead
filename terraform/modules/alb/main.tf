resource "aws_lb" "main" {
  name               = "musclead-alb"
  internal           = false
  load_balancer_type = "application" # http/https用のALB
  security_groups    = [var.alb_sg_id]
  subnets            = var.subnet_ids

  enable_deletion_protection = false

  tags = {
    Name = "musclead-alb"
  }
}
