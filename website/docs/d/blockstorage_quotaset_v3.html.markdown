---
layout: "openstack"
page_title: "OpenStack: openstack_blockstorage_quotaset_v3"
sidebar_current: "docs-openstack-datasource-blockstorage-quotaset-v3"
description: |-
  Get information on a BlockStorage Quotaset v3 of a project.
---

# openstack\_blockstorage\_quotaset\_v3

Use this data source to get the blockstorage quotaset v3 of an OpenStack project.

## Example Usage

```hcl
data "openstack_blockstorage_quotaset_v3" "quota" {
  project_id = "2e367a3d29f94fd988e6ec54e305ec9d"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V3 Blockstorage client.
    If omitted, the `region` argument of the provider is used.

* `project_id` - (Required) The id of the project to retrieve the quotaset.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `volumes` -  The number of volumes that are allowed.
* `snapshots` - The number of snapshots that are allowed.
* `gigabytes` - The size (GB) of volumes and snapshots that are allowed.
* `per_volume_gigabytes` - The size (GB) of volumes that are allowed for each volume.
* `backups` - The number of backups that are allowed.
* `backup_gigabytes` - The size (GB) of backups that are allowed.
* `groups` - The number of groups that are allowed.
* `volume_type_quota` - Map with gigabytes_{volume_type}, snapshots_{volume_type}, volumes_{volume_type} for each volume type.