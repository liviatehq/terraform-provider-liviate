---
layout: "liviate"
page_title: "CloudStack: liviate_autoscale_policy"
sidebar_current: "docs-liviate-data-source-autoscale-policy"
description: |-
  Gets information about a CloudStack autoscale policy.
---

# liviate_autoscale_policy

Use this data source to get information about a CloudStack autoscale policy.

## Example Usage

```hcl
# Get policy by ID
data "liviate_autoscale_policy" "existing_policy" {
  id = "6a8dc025-d7c9-4676-8a7d-e2d9b55e7e60"
}

# Get policy by name
data "liviate_autoscale_policy" "scale_up_policy" {
  filter {
    name  = "name"
    value = "scale-up-policy"
  }
}

# Use in an autoscale VM group
resource "liviate_autoscale_vm_group" "vm_group" {
  name             = "web-autoscale"
  lbrule_id        = liviate_loadbalancer_rule.lb.id
  min_members      = 1
  max_members      = 5
  vm_profile_id    = liviate_autoscale_vm_profile.profile.id
  
  scaleup_policy_ids = [
    data.liviate_autoscale_policy.existing_policy.id
  ]
  
  scaledown_policy_ids = [
    liviate_autoscale_policy.scale_down.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the autoscale policy.

* `filter` - (Optional) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `id` - The autoscale policy ID.

* `name` - The name of the policy.

* `action` - The action (SCALEUP or SCALEDOWN).

* `duration` - The duration in seconds.

* `quiet_time` - The quiet time in seconds.

* `condition_ids` - The list of condition IDs used by this policy.
