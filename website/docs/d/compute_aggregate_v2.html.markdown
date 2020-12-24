---
layout: "openstack"
page_title: "OpenStack: openstack_compute_aggregate_v2"
sidebar_current: "docs-openstack-datasource-compute-aggregate-v2"
description: |-
  Get information on Openstack Host Aggregate
---

# openstack\_compute\_aggregate\_v2

Use this data source to get information about host aggregates
by name.

## Example Usage

```hcl
data "openstack_compute_aggregate_v2" "test" {
  name = "test"
}
```

## Argument Reference

* `name` - The name of the host aggregate

## Attributes Reference

`id` is set to the ID of the found Host Aggregate. In addition, the
following attributes are exported:

* `name` - See Argument Reference above.
* `zone` - Availability zone of the Host Aggregate
* `metadata` - Metadata of the Host Aggregate
* `hosts` - List of Hypervisors contained in the Host Aggregate
