---
layout: "openstack"
page_title: "OpenStack: openstack_sharedfilesystem_share_v2"
sidebar_current: "docs-openstack-datasource-sharedfilesystem-share-v2"
description: |-
  Get information on an Shared File System share.
---

# openstack\_sharedfilesystem\_share\_v2

Use this data source to get the ID of an available Shared File System share.

## Example Usage

```hcl
data "openstack_sharedfilesystem_share_v2" "share_1" {
  name = "share_1"
}
```

## Argument Reference

* `name` - (Optional) The name of the share.

* `description` - (Optional) The human-readable description for the share.

* `project_id` - (Optional) The owner of the share.

* `status` - (Optional) A share status filter. A valid value is `creating`,
   `error`, `available`, `deleting`, `error_deleting`, `manage_starting`,
   `manage_error`, `unmanage_starting`, `unmanage_error`, `unmanaged`,
   `extending`, `extending_error`, `shrinking`, `shrinking_error`, or
   `shrinking_possible_data_loss_error`.

* `is_public` - (Optional) The level of visibility for the share.
    length.

## Attributes Reference

`id` is set to the ID of the found share. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `status` - See Argument Reference above.
* `is_public` - See Argument Reference above.
* `region` - The region in which to obtain the V2 Shared File System client.
* `availability_zone` - The share availability zone.
* `share_proto` - The share protocol.
* `size` - The share size, in GBs.
* `metadata` - One or more metadata key and value pairs as a dictionary of
    strings.
* `export_locations` - A list of export locations. For example, when a share
    server has more than one network interface, it can have multiple export
    locations.
