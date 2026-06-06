output "bucket_name" {
  description = "S3 bucket 名(BE で env として参照)"
  value       = aws_s3_bucket.images.id
}

output "bucket_arn" {
  description = "S3 bucket ARN(IAM policy で参照)"
  value       = aws_s3_bucket.images.arn
}
