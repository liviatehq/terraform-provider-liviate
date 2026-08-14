# Sample: Kubernetes

Provisions a Kubernetes cluster using the Liviate provider.

## Files

| File | Purpose |
|------|---------|
| `main.tf` | Provider + `required_providers` block |
| `variables.tf` | Input variables (API credentials are `sensitive`) |
| `kubernetes.tf` | Data source + a sample `liviate_kubernetes_cluster` |
| `terraform.tfvars.example` | Copy to `terraform.tfvars` and set real values |

## What it does

- Looks up the default zone with `data "liviate_zone"`.
- Creates a `liviate_kubernetes_cluster` with:
  - Kubernetes version `v1.36.0` (by name — resolved via the API)
  - Service offering `liviate-vcpu-2-4gb-20260317-01` (by name)
  - 3 worker nodes

## Run it

```sh
cd sample/kubernetes
cp terraform.tfvars.example terraform.tfvars   # then fill in your API key/secret
terraform init
terraform plan
terraform apply
```

## Notes

- `kubernetes_version` and `service_offering` accept **names** (or IDs); the
  provider resolves them against the API.
- The Kubernetes version name may be with or without a leading `v` — check what
  is registered in your environment and adjust `kubernetes_version` in
  `terraform.tfvars` if the plan errors.
