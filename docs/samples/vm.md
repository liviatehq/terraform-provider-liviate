# Sample: VM

Provisions a virtual machine using the Liviate provider.

## Files

| File | Purpose |
|------|---------|
| `main.tf` | Provider + `required_providers` block |
| `variables.tf` | Input variables (API credentials are `sensitive`) |
| `instances.tf` | Data sources + a sample `liviate_instance` |
| `terraform.tfvars.example` | Copy to `terraform.tfvars` and set real values |

## What it does

- Looks up a template by name (`debian-13-std-30`) with `data "liviate_template"`.
- Looks up a service offering (`liviate-vcpu-2-4gb-20260317-01`) with
  `data "liviate_service_offering"`.
- Looks up the default zone with `data "liviate_zone"`.
- Creates a VM with a 20 GB root disk.

## Run it

```sh
cd sample/vm
cp terraform.tfvars.example terraform.tfvars   # then fill in your API key/secret
terraform init
terraform plan
terraform apply
```

## Notes

- Data source filters match by regular expression on the `name` field.
  Adjust the `value` in `instances.tf` to match templates/offerings available
  in your environment.
- Setting `expunge = true` on the instance permanently deletes it on
  `terraform destroy` instead of leaving it in the CloudStack "Destroyed" state.
- `root_disk_size` can be changed in place (the provider stops, resizes, and
  starts the VM via the `resizeVolume` API).
