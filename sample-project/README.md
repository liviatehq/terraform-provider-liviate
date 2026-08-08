# Liviate Sample Project

Example Terraform configuration using the locally-built Liviate provider.

## Files

| File | Purpose |
|------|---------|
| `main.tf` | Provider + required_providers block (source `liviate/liviate`) |
| `variables.tf` | Input variables (API credentials are `sensitive`) |
| `instances.tf` | Data sources + a sample `liviate_instance` |
| `terraform.tfvars.example` | Copy to `terraform.tfvars` and set real values |

## Usage

1. Make sure the provider is installed in the local plugin mirror
   (already done — see `%APPDATA%\terraform.d\plugins\registry.terraform.io\liviate\liviate\0.5.0\windows_amd64`).

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
