---
subcategory: "Shared Filesystem / Manila"
layout: "openstack"
page_title: "OpenStack: sharedfilesystem_sharetype_v2"
sidebar_current: "docs-openstack-resource-sharedfilesystem-sharetype-v2"
description: |-
  Configure a Shared File System share type.
---

# sharedfilesystem\_sharetype\_v2

Use this resource to configure a Shared File System (Manila) share type.

A share type is analogous to a Block Storage volume type: it groups a set of
capabilities (the "extra specs") that a share must satisfy, and is used to
filter and weigh the back ends where new shares are created. Creating and
managing share types is an admin-only operation.

Minimum supported Manila microversion is 2.7. Reading and updating the
`description` field, and using the `openstack_sharedfilesystem_sharetype_v2`
resource's `terraform import` support to reveal it, require microversion
2.50 or later.

## Example Usage

```hcl
resource "openstack_sharedfilesystem_sharetype_v2" "sharetype_1" {
  name        = "my_share_type"
  description = "created by terraform"
  is_public   = true

  extra_specs = {
    driver_handles_share_servers = "true"
    snapshot_support             = "true"
  }
}
```

### With additional backend-specific extra specs

```hcl
resource "openstack_sharedfilesystem_sharetype_v2" "sharetype_1" {
  name = "my_share_type"

  extra_specs = {
    driver_handles_share_servers = "false"
    snapshot_support             = "true"
    replication_type             = "readable"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Shared File
    System client. A Shared File System client is needed to create a share
    type. If omitted, the `region` argument of the provider is used.
    Changing this creates a new share type.

* `name` - (Required) The name of the share type. Changing this updates the
    name of the existing share type.

* `description` - (Optional) The human-readable description for the share
    type. Changing this updates the description of the existing share type.
    Requires Manila microversion 2.50 or later.

* `is_public` - (Optional) Indicates whether the share type is publicly
    accessible. Defaults to `true`. Changing this updates the existing share
    type.

* `extra_specs` - (Required) A map of extra specifications that filter
    which back ends the share type applies to. Must always contain a
    `driver_handles_share_servers` key with a value of `"true"` or
    `"false"`. May also contain a `snapshot_support` key, as well as any
    other backend-specific extra specification. Changing this updates the
    existing share type's extra specifications.

## Attributes Reference

* `id` - The unique ID for the share type.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `is_public` - See Argument Reference above.
* `extra_specs` - See Argument Reference above.
* `is_default` - Whether this share type is the Manila default share type.

## Import

This resource can be imported by specifying the ID of the share type:

```
$ terraform import openstack_sharedfilesystem_sharetype_v2.sharetype_1 id
```
