---
layout: default
title: "CloudStack: liviate_user"
sidebar_current: "docs-liviate-resource-user"
description: |-
    Creates a User
---

# CloudStack: liviate_user

A `liviate_user` resource manages a user within CloudStack.

## Example Usage

```hcl
resource "liviate_user" "example" {
    account = "example-account"
    email = "user@example.com"
    first_name = "John"
    last_name = "Doe"
    password = "securepassword"
    username = "jdoe"
}
```


## Argument Reference

The following arguments are supported:

* `account` - (Optional) The account the user belongs to.
* `email` - (Required) The email address of the user.
* `first_name` - (Required) The first name of the user.
* `last_name` - (Required) The last name of the user.
* `password` - (Required) The password for the user.
* `username` - (Required) The username of the user.

## Attributes Reference

No attributes are exported.

## Import

Users can be imported; use `<USERID>` as the import ID. For example:

```shell
$ terraform import liviate_user.example <USERID>
```
