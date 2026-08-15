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

variable "zone_id" {
  description = "Zone ID to deploy into"
  type        = string
  default     = ""
}

variable "ssh_key_name" {
  description = "SSH keypair name"
  type        = string
  default     = "liviate-key"
}

variable "public_key_path" {
  description = "Path to the public key to register"
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}
