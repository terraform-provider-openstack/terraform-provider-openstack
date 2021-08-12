---
layout: "openstack"
page_title: "OpenStack: openstack_blockstorage_volume_type_v3"
sidebar_current: "docs-openstack-resource-blockstorage-volume-type-v3"
description: |-
  Manages a V3 volume type resource within OpenStack.
---

# openstack\_blockstorage\_volume\_type\_v3

Manages a V3 block storage volume type resource within OpenStack.

~> **Note:** This usually requires admin privileges.


## Example Usage

```hcl
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name        = "volume_type_1"
  description = "Volume type 1"
  extra_specs = {
      capabilities        = "gpu"
      volume_backend_name = "ssd"
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the volume. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new quotaset.

* `name` - (Required) Name of the volume type.  Changing this
    updates the `name` of an existing volume type.

* `description` - (Optional) Human-readable description of the port. Changing
    this updates the `description` of an existing volume type.

* `is_public` - (Optional) Whether the volume type is public. Changing
    this updates the `is_public` of an existing volume type.

* `extra_specs` - (Optional) Key/Value pairs of metadata for the volume type.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `is_public` - See Argument Reference above.
* `extra_specs` - See Argument Reference above.

## Import

Volume types can be imported using the `volume_type_id`, e.g.

```
$ terraform import openstack_blockstorage_volume_type_v3.volume_type_1 941793f0-0a34-4bc4-b72e-a6326ae58283
```
