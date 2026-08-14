---
layout: "liviate"
page_title: "Liviate: liviate_volume"
sidebar_current: "docs-liviate-liviate_volume"
description: |-
  Gets information about liviate volume.
---

# liviate_volume

Use this datasource to get information about a volume for use in other resources.

### Example Usage

```hcl
  data "liviate_volume" "volume-data-source"{
    filter{
    name = "name"
    value="TestVolume"
    }
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the disk volume.
* `disk_offering_id` - ID of the disk offering.
* `zone_id` - ID of the availability zone.
