---
layout: "openstack"
page_title: "OpenStack: openstack_compute_availability_zones_v2"
sidebar_current: "docs-openstack-datasource-compute-availability-zones-v2"
description: |-
  Get a list of availability zones from OpenStack
---

# openstack\_compute\_availability\_zones\_v2

Use this data source to get a list of availability zones from OpenStack

## Example Usage

```hcl
data "openstack_compute_availability_zones_v2" "zones" {}
```

## Argument Reference

* `region` - (Optional) The `region` to fetch availability zones from, defaults to the provider's `region`
* `state` - (Optional) The `state` of the availability zones to match, default ("available").


## Attributes Reference

`id` is set to hash of the returned zone list. In addition, the following attributes
are exported:

* `names` - The names of the availability zones, ordered alphanumerically, that match the queried `state`
