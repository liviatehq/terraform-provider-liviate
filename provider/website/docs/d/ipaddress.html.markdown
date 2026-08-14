---
layout: "liviate"
page_title: "Liviate: liviate_ipaddress"
sidebar_current: "docs-liviate-liviate_ipaddress"
description: |-
  Gets information about liviate ipaddress.
---

# liviate_ipaddress

Use this datasource to get information about a public ipaddress for use in other resources.

### Example Usage

```hcl
data "liviate_ipaddress" "ipaddress-data-source"{
    filter{
    name = "zone_name"
    value= "DC"
    }
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `is_portable` - Is public IP portable across the zones.
* `network_id` - The ID of the Network where ip belongs to.
* `vpc_id` - VPC id the ip belongs to.
* `zone_name` - The name of the zone the public IP address belongs to.
* `project` - The project name of the address.
* `ip_address` - Public IP address.
* `is_source_nat` - True if the IP address is a source nat address, false otherwise.
