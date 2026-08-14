---
layout: "liviate"
page_title: "Liviate: liviate_user"
sidebar_current: "docs-liviate-liviate_user"
description: |-
  Gets information about liviate user.
---

# liviate_user

Use this datasource to get information about a liviate user for use in other resources.

### Example Usage

```hcl
data "liviate_user" "user-data-source"{
    filter{
    name = "first_name"
    value= "jon"
    }
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `account` - The account name of the userg.
* `email` - The user email address.
* `first_name` - The user firstname.
* `last_name` - The user lastname.
* `username` - The user name