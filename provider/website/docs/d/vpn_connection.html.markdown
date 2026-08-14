---
layout: "liviate"
page_title: "Liviate: liviate_vpn_connection"
sidebar_current: "docs-liviate-liviate_vpn_connection"
description: |-
  Gets information about liviate vpn connection.
---

# liviate_vpn_connection

Use this datasource to get information about a vpn connection for use in other resources.

### Example Usage

```hcl
data "liviate_vpc" "vpc-data-source"{
    filter{
    name = "name"
    value= "test-vpc"
    }
    filter{
    name = "cidr"
    value= "10.0.0.0/16"
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `s2s_customer_gateway_id` - The customer gateway ID.
* `s2s_vpn_gateway_id` - The vpn gateway ID.