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
}

variable "kubernetes_version" {
  description = "Kubernetes version name or ID (e.g. v1.36.0)"
  type        = string
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "liviate-k8s-csi-01"
}

variable "service_offering" {
  description = "Service offering name for cluster nodes"
  type        = string
  default     = "liviate-vcpu-2-4gb-20260317-01"
}

variable "disk_offering_id" {
  description = "A CUSTOMIZED-size disk offering ID -- the CSI driver needs iscustomized=true so it can honor whatever size a PersistentVolumeClaim actually requests, instead of snapping to a fixed offering size."
  type        = string
}
