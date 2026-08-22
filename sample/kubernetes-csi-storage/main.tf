terraform {
  required_version = ">= 1.0"
  required_providers {
    liviate = {
      source  = "liviatehq/liviate"
      version = ">= 1.1.0" # enable_csi / liviate_role_permission / liviate_kubernetes_cluster_config
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.30"
    }
  }
}

provider "liviate" {
  api_key    = var.liviate_api_key
  secret_key = var.liviate_secret_key
}
