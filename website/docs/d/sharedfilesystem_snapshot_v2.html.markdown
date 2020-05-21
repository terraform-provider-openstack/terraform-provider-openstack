---
layout: "openstack"
page_title: "OpenStack: openstack_sharedfilesystem_snapshot_v2"
sidebar_current: "docs-openstack-datasource-sharedfilesystem-snapshot-v2"
description: |-
  Get information on an Shared File System snapshot.
---

# openstack\_sharedfilesystem\_snapshot\_v2

Use this data source to get the ID of an available Shared File System snapshot.

## Example Usage

```hcl
data "openstack_sharedfilesystem_snapshot_v2" "snapshot_1" {
  name = "snapshot_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Shared File System client.

* `name` - (Optional) The name of the snapshot.

* `description` - (Optional) The human-readable description of the snapshot.

* `project_id` - (Optional) The owner of the snapshot.

* `status` - (Optional) A snapshot status filter. A valid value is `available`, `error`,
    `creating`, `deleting`, `manage_starting`, `manage_error`, `unmanage_starting`,
    `unmanage_error` or `error_deleting`.

## Attributes Reference

`id` is set to the ID of the found snapshot. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `status` - See Argument Reference above.
* `size` - The snapshot size, in GBs.
* `share_id` - The UUID of the source share that was used to create the snapshot.
* `share_proto` - The file system protocol of a share snapshot.
* `share_size` - The share snapshot size, in GBs.
