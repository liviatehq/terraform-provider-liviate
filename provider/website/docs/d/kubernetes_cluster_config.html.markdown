---
layout: default
page_title: "CloudStack: liviate_kubernetes_cluster_config"
sidebar_current: "docs-liviate-datasource-kubernetes_cluster_config"
description: |-
    Retrieves a CKS cluster's kubeconfig, parsed into fields ready for the kubernetes/helm providers
---

# CloudStack: liviate_kubernetes_cluster_config

Retrieves a `liviate_kubernetes_cluster`'s kubeconfig (`getKubernetesClusterConfig`) and parses
out the fields the standard `hashicorp/kubernetes` and `hashicorp/helm` providers expect, so a
Terraform config can chain straight from provisioning a cluster to managing resources inside it --
no manual `kubectl`/copy-paste step.

CKS always returns a single-cluster, single-user, single-context kubeconfig using client-certificate
auth, which is what this data source assumes.

## Example Usage

```hcl
resource "liviate_kubernetes_cluster" "cluster" {
  name               = "my-cluster"
  zone               = "zone1"
  kubernetes_version = "v1.36.0"
  service_offering   = "liviate-vcpu-2-4gb-20260317-01"
  enable_csi         = true
}

data "liviate_kubernetes_cluster_config" "cluster" {
  cluster_id = liviate_kubernetes_cluster.cluster.id
}

provider "kubernetes" {
  host                   = data.liviate_kubernetes_cluster_config.cluster.host
  cluster_ca_certificate = base64decode(data.liviate_kubernetes_cluster_config.cluster.cluster_ca_certificate)
  client_certificate     = base64decode(data.liviate_kubernetes_cluster_config.cluster.client_certificate)
  client_key             = base64decode(data.liviate_kubernetes_cluster_config.cluster.client_key)
}
```

## Argument Reference

* `cluster_id` - (Required) ID of the Kubernetes cluster.

## Attributes Reference

* `raw_config` - The full kubeconfig YAML, unparsed. Sensitive.
* `host` - The Kubernetes API server URL.
* `cluster_ca_certificate` - Base64-encoded CA certificate. Pass through `base64decode()` into a
  provider's `cluster_ca_certificate` argument. Sensitive.
* `client_certificate` - Base64-encoded client certificate. Sensitive.
* `client_key` - Base64-encoded client key. Sensitive.
