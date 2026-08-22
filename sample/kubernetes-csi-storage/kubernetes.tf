# 1. The cluster itself, with the CSI driver installed at bootstrap (enable_csi is ForceNew --
#    CloudStack has no API to toggle it on a running cluster, only at creation).
resource "liviate_kubernetes_cluster" "cluster" {
  name               = var.cluster_name
  zone               = var.zone_id
  kubernetes_version = var.kubernetes_version
  service_offering   = var.service_offering
  size               = 1
  control_nodes_size = 1
  enable_csi         = true
}

# 2. Every CKS cluster's CSI/CCM sidecar authenticates as a CloudStack account under the built-in
#    "Project Kubernetes Service Role" -- which, as shipped, has no volume-management permissions
#    at all (it predates CSI support). Without this, every PVC fails to provision with "API
#    [listVolumes] does not exist or is not available for the account". Look the role up by name
#    (it's a platform default, not something Terraform creates) and grant what CSI actually needs.
data "liviate_role" "kubernetes_service" {
  filter {
    name  = "name"
    value = "^Project Kubernetes Service Role$"
  }
}

locals {
  # listVolumes must exist for the reorder-before-wildcard step to have something to reorder
  # around in the first place, so it's listed even though it'd also be covered by a future "list*".
  csi_role_rules = [
    "listVolumes", "createVolume", "deleteVolume", "attachVolume", "detachVolume",
    "resizeVolume", "listDiskOfferings", "listZones", "listServiceOfferings",
  ]
}

resource "liviate_role_permission" "csi_volume_access" {
  for_each = toset(local.csi_role_rules)

  role_id     = data.liviate_role.kubernetes_service.id
  rule        = each.value
  permission  = "allow"
  description = "CSI driver support"
}

# 3. Pull the cluster's own kubeconfig (client-cert auth) to configure the standard Kubernetes
#    provider against it -- no manual `getKubernetesClusterConfig` + kubectl step needed.
data "liviate_kubernetes_cluster_config" "cluster" {
  cluster_id = liviate_kubernetes_cluster.cluster.id
}

provider "kubernetes" {
  host                   = data.liviate_kubernetes_cluster_config.cluster.host
  cluster_ca_certificate = base64decode(data.liviate_kubernetes_cluster_config.cluster.cluster_ca_certificate)
  client_certificate     = base64decode(data.liviate_kubernetes_cluster_config.cluster.client_certificate)
  client_key             = base64decode(data.liviate_kubernetes_cluster_config.cluster.client_key)
}

# 4. The actual point of all the above: a working default StorageClass so a plain
#    `kubectl apply` of a PVC just works. `csi.cloudstack.apache.org/disk-offering-id` is the
#    ONLY parameter the driver's provisioner will accept (found by trial -- a bare "diskoffering"
#    key is silently ignored and fails with "Missing parameter csi.cloudstack.apache.org/disk-
#    offering-id"). This depends on both the cluster AND the role-permission grant above --
#    without an explicit depends_on, Terraform has no way to know the CSI driver pods (which
#    read those permissions at PVC-provision time, not at StorageClass-creation time) need the
#    grant to exist first.
resource "kubernetes_storage_class" "cloudstack_csi" {
  metadata {
    name = "cloudstack-csi"
    annotations = {
      "storageclass.kubernetes.io/is-default-class" = "true"
    }
  }
  storage_provisioner = "csi.cloudstack.apache.org"
  reclaim_policy       = "Delete"
  volume_binding_mode  = "Immediate"
  allow_volume_expansion = true
  parameters = {
    "csi.cloudstack.apache.org/disk-offering-id" = var.disk_offering_id
  }

  depends_on = [liviate_role_permission.csi_volume_access]
}
