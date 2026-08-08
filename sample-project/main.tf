terraform {
  required_version = ">= 1.0"
  required_providers {
    liviate = {
      source  = "liviate/liviate"
      version = "0.5.0"
    }
  }
}

provider "liviate" {
  api_key    = var.liviate_api_key
  secret_key = var.liviate_secret_key
}
