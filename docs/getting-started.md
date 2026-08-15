# Getting Started

This guide covers building the provider from source, pointing Terraform at your
local build, and running a sample configuration.

## Requirements

- [Go](https://golang.org/doc/install) 1.23+
- [Terraform](https://www.terraform.io/downloads) 1.x
- A Liviate (CloudStack) account with an API key and secret

## Build the provider

```sh
cd provider
go build -o terraform-provider-liviate.exe .
```

This produces `terraform-provider-liviate.exe` in the `provider/` directory.

## Point Terraform at your local build (dev override)

Terraform normally installs providers from the Terraform registry. For local
development, use a **dev override** so Terraform runs your freshly built binary
directly and skips install/checksum handling.

Create or edit `%APPDATA%\terraform.rc` (Windows) or `~/.terraformrc`
(Linux/macOS):

```hcl
provider_installation {
  dev_overrides {
    "liviatehq/liviate" = "C:/Users/madsk/workspace/liviate-cs/provider"
  }
  direct {}
}
```

Replace the path with the absolute path to the `provider/` directory that
contains `terraform-provider-liviate.exe`.

> **Tip:** with a dev override in place you do **not** need to run
> `terraform init` after every rebuild — each `terraform plan`/`apply` uses the
> latest binary. You may see a warning about development overrides; that is
> expected and harmless.

## Run a sample

```sh
cd sample/vm
terraform init
terraform plan
terraform apply
```

Credentials come from `terraform.tfvars` (gitignored — copy
`terraform.tfvars.example` and fill in your API key and secret).

## Rebuilding after provider changes

```sh
cd provider
go build -o terraform-provider-liviate.exe .
cd ../sample/vm
terraform plan   # picks up the new binary automatically
```
