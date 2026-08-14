# Liviate Terraform Provider

The Liviate Terraform provider manages cloud resources on the Liviate platform
(backed by Apache CloudStack). It is a fork of the Apache CloudStack Terraform
provider, rebranded for Liviate.

This repository contains:

| Path | Description |
|------|-------------|
| [`provider/`](provider/) | The Liviate Terraform provider (Go source, build config, registry docs) |
| [`sample/`](sample/) | Example Terraform configurations using the provider |
| [`docs/`](docs/) | Documentation: getting started, provider reference, samples |

## Quick start

1. Build the provider and set up the dev override (see [docs/getting-started.md](docs/getting-started.md)).
2. Run one of the samples, e.g. [sample/vm](sample/vm/).
3. For the full resource reference, see [docs/resources.md](docs/resources.md).

## Documentation

- [Getting started](docs/getting-started.md)
- [Provider configuration](docs/provider.md)
- [Resources & data sources](docs/resources.md)
- [Samples](docs/samples/)

## License

Licensed under the Apache License, Version 2.0. See [provider/LICENSE](provider/LICENSE).
