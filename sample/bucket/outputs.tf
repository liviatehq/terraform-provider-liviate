output "private_bucket_name" {
  value = liviate_bucket.private.name
}

output "private_bucket_url" {
  value = liviate_bucket.private.url
}

output "private_bucket_access_key" {
  value = liviate_bucket.private.accesskey
}

output "private_bucket_secret_key" {
  value     = liviate_bucket.private.usersecretkey
  sensitive = true
}

output "public_bucket_name" {
  value = liviate_bucket.public.name
}

output "public_bucket_url" {
  value = liviate_bucket.public.url
}

output "public_bucket_access_key" {
  value = liviate_bucket.public.accesskey
}

output "public_bucket_secret_key" {
  value     = liviate_bucket.public.usersecretkey
  sensitive = true
}
