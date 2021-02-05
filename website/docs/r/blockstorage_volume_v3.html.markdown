---
layout: "openstack"
page_title: "OpenStack: openstack_blockstorage_volume_v3"
sidebar_current: "docs-openstack-resource-blockstorage-volume-v3"
description: |-
  Manages a V3 volume resource within OpenStack.
---

# openstack\_blockstorage\_volume\_v3

Manages a V3 volume resource within OpenStack.

## Example Usage

```hcl
resource "openstack_blockstorage_volume_v3" "volume_1" {
  region      = "RegionOne"
  name        = "volume_1"
  description = "first test volume"
  size        = 3
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the volume. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new volume.

* `size` - (Required) The size of the volume to create (in gigabytes).

* `enable_online_resize` - (Optional) When this option is set it allows extending
    attached volumes. Note: updating size of an attached volume requires Cinder
    support for version 3.42 and a compatible storage driver.

* `availability_zone` - (Optional) The availability zone for the volume.
    Changing this creates a new volume.

* `consistency_group_id` - (Optional) The consistency group to place the volume
    in.

* `description` - (Optional) A description of the volume. Changing this updates
    the volume's description.

* `image_id` - (Optional) The image ID from which to create the volume.
    Changing this creates a new volume.

* `metadata` - (Optional) Metadata key/value pairs to associate with the volume.
    Changing this updates the existing volume metadata.

* `name` - (Optional) A unique name for the volume. Changing this updates the
    volume's name.

* `snapshot_id` - (Optional) The snapshot ID from which to create the volume.
    Changing this creates a new volume.

* `source_replica` - (Optional) The volume ID to replicate with.

* `source_vol_id` - (Optional) The volume ID from which to create the volume.
    Changing this creates a new volume.

* `volume_type` - (Optional) The type of volume to create.
    Changing this creates a new volume.

* `multiattach` - (Optional) Allow the volume to be attached to more than one Compute instance.

* `scheduler_hints` - (Optional) Provide the Cinder scheduler with hints on where
    to instantiate a volume in the OpenStack cloud. The available hints are described below.
    
The `scheduler_hints` block supports:

* `different_host` - (Optional) The volume should be scheduled on a 
    different host from the set of volumes specified in the list provided.

* `same_host` - (Optional) A list of volume UUIDs. The volume should be
    scheduled on the same host as another volume specified in the list provided.

* `local_to_instance` - (Optional) An instance UUID. The volume should be 
    scheduled on the same host as the instance.

* `query` - (Optional) A conditional query that a back-end must pass in
    order to host a volume. The query must use the `JsonFilter` syntax
    which is described
    [here](https://docs.openstack.org/cinder/latest/configuration/block-storage/scheduler-filters.html#jsonfilter).
    At this time, only simple queries are supported. Compound queries using
    `and`, `or`, or `not` are not supported. An example of a simple query is:

    ```
    [“=”, “$backend_id”, “rbd:vol@ceph#cloud”]
    ```

* `additional_properties` - (Optional) Arbitrary key/value pairs of additional
  properties to pass to the scheduler.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `size` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `image_id` - See Argument Reference above.
* `source_vol_id` - See Argument Reference above.
* `snapshot_id` - See Argument Reference above.
* `metadata` - See Argument Reference above.
* `volume_type` - See Argument Reference above.
* `attachment` - If a volume is attached to an instance, this attribute will
    display the Attachment ID, Instance ID, and the Device as the Instance
    sees it.
* `multiattach` - See Argument Reference above.

## Import

Volumes can be imported using the `id`, e.g.

```
$ terraform import openstack_blockstorage_volume_v3.volume_1 ea257959-eeb1-4c10-8d33-26f0409a755d
```
