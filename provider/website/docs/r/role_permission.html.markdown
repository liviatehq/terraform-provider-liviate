---
layout: default
page_title: "CloudStack: liviate_role_permission"
sidebar_current: "docs-liviate-resource-role_permission"
description: |-
    Manages a single API permission rule on a CloudStack role
---

# CloudStack: liviate_role_permission

A `liviate_role_permission` resource manages one `allow`/`deny` rule on a CloudStack role
(`liviate_role`, or a built-in role looked up via the `liviate_role` data source).

## Rule ordering

CloudStack evaluates a role's rules **in list order, first match wins**. `createRolePermission`
always appends a new rule at the end of the list -- so if the role already has a trailing wildcard
rule (e.g. `*` deny, the normal pattern for a custom role), any rule created afterward is dead on
arrival: the wildcard matches first and the new rule never gets evaluated.

This resource handles that automatically: after creating a rule, if the role has any wildcard
(`*`) rule, it reorders the whole rule list so every specific (non-wildcard) rule sits ahead of
every wildcard rule, preserving relative order within each group. This is a deliberate
simplification rather than a general-purpose ordering system -- it solves the one case that
actually comes up in practice (add a specific override on top of an existing catch-all) without
requiring every user of this resource to reason about full list ordering.

If you need more granular control over ordering among several specific rules, create them in the
order you want and be aware later applies re-run the same "specific-before-wildcard" reorder --
it does not otherwise reshuffle rules relative to each other.

## Example Usage

```hcl
data "liviate_role" "kubernetes_service" {
  filter {
    name  = "name"
    value = "^Project Kubernetes Service Role$"
  }
}

resource "liviate_role_permission" "list_volumes" {
  role_id     = data.liviate_role.kubernetes_service.id
  rule        = "listVolumes"
  permission  = "allow"
  description = "CSI driver support"
}
```

See the `kubernetes-csi-storage` sample for the full working example this is drawn from (a CSI
driver's service account needing volume-management permissions the built-in role doesn't ship
with).

## Argument Reference

* `role_id` - (Required, Forces new resource) ID of the role this permission belongs to.
* `rule` - (Required, Forces new resource) The API name (e.g. `listVolumes`) or a wildcard rule
  (e.g. `list*`, `*`).
* `permission` - (Required) `allow` or `deny`. Unlike `role_id`/`rule`, changing this updates the
  existing rule in place rather than replacing it.
* `description` - (Optional) A description for the rule.

## Attributes Reference

* `id` - The ID of the role permission rule.

## Import

```shell
$ terraform import liviate_role_permission.example <ROLEPERMISSIONID>
```
