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
  description = "Zone ID to deploy the cluster into"
  type        = string
  default     = ""
}

variable "kubernetes_version" {
  description = "Semantic version of Kubernetes (e.g. 1.28.4)"
  type        = string
  default     = ""
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "liviate-k8s-01"
}

variable "node_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 3
}

variable "service_offering" {
  description = "Service offering name for cluster nodes"
  type        = string
  default     = "liviate-vcpu-2-4gb-20260317-01"
}
