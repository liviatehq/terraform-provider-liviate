---
layout: default
page_title: "CloudStack: liviate_zone"
sidebar_current: "docs-liviate-resource-zone"
description: |-
    Creates a Zone
---

# CloudStack: liviate_zone

A `liviate_zone` resource manages a zone within CloudStack.

## Example Usage
```hcl
resource "liviate_zone" "example" {
    name = "example-zone"
    dns1 = "8.8.8.8"
    internal_dns1 = "8.8.4.4"
    network_type = "Basic"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the zone.
* `dns1` - (Required) The DNS server  1 for the zone.
* `internal_dns1` - (Required) The internal DNS server  1 for the zone.
* `network_type` - (Required) The type of network to use for the zone.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the zone.
* `name` - The name of the zone.
* `dns1` - The DNS server  1 for the zone.
* `internal_dns1` - The internal DNS server  1 for the zone.
* `network_type` - The type of network to use for the zone.

## Import

Zones can be imported; use `<ZONEID>` as the import ID. For example:
```shell
$ terraform import liviate_zone.example <ZONEID>
```
