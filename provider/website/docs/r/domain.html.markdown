---
layout: default
page_title: "CloudStack: liviate_domain"
sidebar_current: "docs-liviate-resource-domain"
description: |-
    Creates a Domain
---

# CloudStack: liviate_domain

A `liviate_domain` resource manages a domain within CloudStack.

## Example Usage

```hcl
resource "liviate_domain" "example" {
    name = "example-domain"
    network_domain = "example.local"
    parent_domain_id = "ROOT"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the domain.
* `domain_id` - (Optional) The ID of the domain.
* `network_domain` - (Optional) The network domain for the domain.
* `parent_domain_id` - (Optional) The ID of the parent domain.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the domain.
* `name` - The name of the domain.
* `network_domain` - The network domain for the domain.
* `parent_domain_id` - The ID of the parent domain.

## Import

Domains can be imported; use `<DOMAINID>` as the import ID. For example:

```shell
$ terraform import liviate_domain.example <DOMAINID>
```
