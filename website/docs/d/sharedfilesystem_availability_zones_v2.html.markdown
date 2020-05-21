---
layout: "openstack"
page_title: "OpenStack: openstack_sharedfilesystem_availability_zones_v2"
sidebar_current: "docs-openstack-datasource-sharedfilesystem-availability-zones-v2"
description: |-
  Get a list of Shared File System availability zones from OpenStack
---

# openstack\_sharedfilesystem\_availability\_zones\_v2

Use this data source to get a list of Shared File System availability zones
from OpenStack

## Example Usage

```hcl
data "openstack_sharedfilesystem_availability_zones_v2" "zones" {}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Shared File System
    client. If omitted, the `region` argument of the provider is used.

## Attributes Reference

`id` is set to hash of the returned zone list. In addition, the following
attributes are exported:

* `region` - See Argument Reference above.
* `names` - The names of the availability zones, ordered alphanumerically.
