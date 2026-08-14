resource "liviate_bucket" "private" {
  name  = var.bucket_name
  quota = var.bucket_quota
}
