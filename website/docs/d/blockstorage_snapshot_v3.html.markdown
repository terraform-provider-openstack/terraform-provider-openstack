---
layout: "openstack"
page_title: "OpenStack: openstack_blockstorage_snapshot_v3"
sidebar_current: "docs-openstack-datasource-blockstorage-snapshot-v3"
description: |-
  Get information on an OpenStack Snapshot.
---

# openstack\_blockstorage\_snapshot\_v3

Use this data source to get information about an existing snapshot.

## Example Usage

```hcl
data "openstack_blockstorage_snapshot_v3" "snapshot_1" {
  name        = "snapshot_1"
  most_recent = true
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V3 Block Storage
    client. If omitted, the `region` argument of the provider is used.

* `name` - (Optional) The name of the snapshot.

* `status` - (Optional) The status of the snapshot.

* `volume_id` - (Optional) The ID of the snapshot's volume.

* `most_recent` - (Optional) Pick the most recently created snapshot if there
    are multiple results.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `status` - See Argument Reference above.
* `volume_id` - See Argument Reference above.
* `description` - The snapshot's description.
* `size` - The size of the snapshot.
* `metadata` - The snapshot's metadata.
