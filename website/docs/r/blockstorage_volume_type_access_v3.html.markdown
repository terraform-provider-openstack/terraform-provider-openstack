---
layout: "openstack"
page_title: "OpenStack: openstack_blockstorage_volume_type_access_v3"
sidebar_current: "docs-openstack-resource-blockstorage-volume-type-access-v3"
description: |-
  Manages a V3 volume type access resource within OpenStack.
---

# openstack\_blockstorage\_volume\_type\_access\_v3

Manages a V3 block storage volume type access resource within OpenStack.

~> **Note:** This usually requires admin privileges.


## Example Usage

```hcl
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name      = "volume_type_1"
  is_public = false
}

resource "openstack_blockstorage_volume_type_access_v3" "volume_type_access" {
  project_id     = "${openstack_identity_project_v3.project_1.id}"
  volume_type_id = "${openstack_blockstorage_volume_type_v3.volume_type_1.id}"
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the volume. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new quotaset.

* `project_id` - (Required) ID of the project to give access to. Changing this
    creates a new resource.

* `volume_type_id` - (Required) ID of the volume type to give access to. Changing
    this creates a new resource.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `volume_type_id` - See Argument Reference above.

## Import

Volume types access can be imported using the `volume_type_id/project_id`, e.g.

```
$ terraform import openstack_blockstorage_volume_type_access_v3.volume_type_access 941793f0-0a34-4bc4-b72e-a6326ae58283/ed498e81f0cc448bae0ad4f8f21bf67f
```
