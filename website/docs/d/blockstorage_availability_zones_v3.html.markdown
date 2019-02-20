---
layout: "openstack"
page_title: "OpenStack: openstack_blockstorage_availability_zones_v3"
sidebar_current: "docs-openstack-datasource-blockstorage-availability-zones-v3"
description: |-
  Get a list of Block Storage availability zones from OpenStack
---

# openstack\_blockstorage\_availability\_zones\_v3

Use this data source to get a list of Block Storage availability zones from OpenStack

## Example Usage

```hcl
data "openstack_blockstorage_availability_zones_v3" "zones" {}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Block Storage client.
    If omitted, the `region` argument of the provider is used.

* `state` - (Optional) The `state` of the availability zones to match. Can
    either be `available` or `unavailable`. Default is `available`.

## Attributes Reference

`id` is set to hash of the returned zone list. In addition, the following
attributes are exported:

* `region` - See Argument Reference above.
* `state` - See Argument Reference above.
* `names` - The names of the availability zones, ordered alphanumerically, that
    match the queried `state`.
