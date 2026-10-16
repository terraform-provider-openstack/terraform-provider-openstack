---
subcategory: "Shared Filesystem / Manila"
layout: "openstack"
page_title: "OpenStack: sharedfilesystem_sharetype_access_v2"
sidebar_current: "docs-openstack-resource-sharedfilesystem-sharetype-access-v2"
description: |-
  Manages access to a private Shared File System share type within OpenStack.
---

# sharedfilesystem\_sharetype\_access\_v2

Use this resource to grant a project (tenant) access to a private Shared
File System (Manila) share type.

By default a private share type (one created with `is_public = false` on
`openstack_sharedfilesystem_sharetype_v2`) is only usable by the project
that owns it. This resource grants additional projects permission to create
shares using that share type, analogous to
`openstack_blockstorage_volume_type_access_v3` for Block Storage volume
types or `openstack_compute_flavor_access_v2` for Compute flavors.

Creating and managing share type access is an admin-only operation.

~> **Note:** This resource has no effect on public share types. Manila
only tracks explicit project access for share types created with
`is_public = false`.

## Example Usage

```hcl
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_sharedfilesystem_sharetype_v2" "sharetype_1" {
  name      = "my_private_share_type"
  is_public = false

  extra_specs = {
    driver_handles_share_servers = "true"
  }
}

resource "openstack_sharedfilesystem_sharetype_access_v2" "access_1" {
  share_type_id = openstack_sharedfilesystem_sharetype_v2.sharetype_1.id
  project_id    = openstack_identity_project_v3.project_1.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Shared File
    System client. A Shared File System client is needed to manage share
    type access. If omitted, the `region` argument of the provider is
    used. Changing this creates a new resource.

* `share_type_id` - (Required) The ID of the share type to grant access
    to. Changing this creates a new resource.

* `project_id` - (Required) The ID of the project to grant access to the
    share type. Changing this creates a new resource.

## Attributes Reference

* `id` - A combination of the share type ID and project ID, separated by
    a slash, e.g. `941793f0-0a34-4bc4-b72e-a6326ae58283/ed498e81f0cc448bae0ad4f8f21bf67f`.
* `region` - See Argument Reference above.
* `share_type_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.

## Import

Share type access can be imported using the `share_type_id/project_id`,
e.g.

```
$ terraform import openstack_sharedfilesystem_sharetype_access_v2.access_1 941793f0-0a34-4bc4-b72e-a6326ae58283/ed498e81f0cc448bae0ad4f8f21bf67f
```
