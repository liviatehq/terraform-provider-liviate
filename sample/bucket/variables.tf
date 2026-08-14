variable "liviate_api_key" {
  description = "Liviate API key"
  type        = string
  sensitive   = true
}

variable "liviate_secret_key" {
  description = "Liviate secret key"
  type        = string
  sensitive   = true
}

variable "bucket_name" {
  description = "Name of the private S3 bucket"
  type        = string
  default     = "liviate-bucket-01"
}

variable "public_bucket_name" {
  description = "Name of the public S3 bucket"
  type        = string
  default     = "liviate-bucket-public-01"
}

variable "bucket_quota" {
  description = "Maximum storage quota in GB"
  type        = number
  default     = 100
}
