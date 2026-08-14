---
layout: "liviate"
page_title: "CloudStack: liviate_quota_enabled"
sidebar_current: "docs-liviate-datasource-quota-enabled"
description: |-
  Checks if quota is enabled in CloudStack.
---

# liviate_quota_enabled

Use this data source to check whether the quota system is enabled in the CloudStack management server.

## Example Usage

```hcl
# Check if quota system is enabled
data "liviate_quota_enabled" "quota_status" {
}

# Use the quota status in conditional logic
resource "liviate_quota_tariff" "cpu_tariff" {
  count = data.liviate_quota_enabled.quota_status.enabled ? 1 : 0
  
  name       = "CPU Usage Tariff"
  usage_type = 1
  value      = 0.05
}
```

## Argument Reference

This data source takes no arguments.

## Attribute Reference

The following attributes are exported:

* `enabled` - A boolean value indicating whether the quota system is enabled in CloudStack.