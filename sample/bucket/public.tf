resource "liviate_bucket" "public" {
  name  = var.public_bucket_name
  quota = var.bucket_quota

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action    = ["s3:GetObject"]
        Resource  = ["arn:aws:s3:::${var.public_bucket_name}/*"]
      }
    ]
  })
}
