output "instance_id" {
  description = "Bastion EC2 の Instance ID(mus-prod alias で参照)"
  value       = aws_instance.bastion.id
}
