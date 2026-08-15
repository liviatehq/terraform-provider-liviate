terraform {
  required_version = ">= 1.0"
  required_providers {
    liviate = {
      source  = "liviatehq/liviate"
      version = ">= 1.0.0"
    }
  }
}

provider "liviate" {
  api_key    = var.liviate_api_key
  secret_key = var.liviate_secret_key
}
