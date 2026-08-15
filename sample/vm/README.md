# Liviate Sample Project

Example Terraform configuration using the locally-built Liviate provider.

## Files

| File | Purpose |
|------|---------|
| `main.tf` | Provider + required_providers block (source `liviatehq/liviate`) |
| `variables.tf` | Input variables (API credentials are `sensitive`) |
| `instances.tf` | Data sources + sample `liviate_instance` resources and an `liviate_ssh_keypair` |
| `terraform.tfvars.example` | Copy to `terraform.tfvars` and set real values |

## Usage

1. The provider is built locally and a dev override is set up
   (see `docs/getting-started.md` for the dev-override setup).

2. Provide credentials. Either copy `terraform.tfvars.example` to `terraform.tfvars`
   and edit it, or set environment variables:

   ```sh
   $env:LIVIATE_API_KEY    = "your-api-key"
   $env:LIVIATE_SECRET_KEY = "your-secret-key"
   ```

   The API URL defaults to `https://console.liviate.com/client/api` inside the
   provider, so only the key and secret are needed.

   (Note: the sample uses `var.liviate_*` inputs; set them via tfvars.)

3. Run:

   ```sh
   terraform init
   terraform plan
   terraform apply
   ```
