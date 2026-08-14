data "liviate_zone" "primary" {
  filter {
    name  = "name"
    value = ".*"
  }
}

resource "liviate_kubernetes_cluster" "cluster" {
  name               = var.cluster_name
  zone               = var.zone_id != "" ? var.zone_id : data.liviate_zone.primary.id
  kubernetes_version = var.kubernetes_version
  service_offering   = var.service_offering
  size               = var.node_count
  description        = "Kubernetes cluster provisioned by Liviate"
}
