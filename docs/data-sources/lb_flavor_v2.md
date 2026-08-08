---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_flavor_v2"
sidebar_current: "docs-openstack-datasource-lb-flavor-v2"
description: |-
  Get information on an OpenStack Load Balancer Flavor.
---

# openstack\_lb\_flavor\_v2

Use this data source to get the ID of an OpenStack Load Balancer flavor.

## Example Usage

```hcl
data "openstack_lb_flavor_v2" "flavor_1" {
  name = "flavor_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Load Balancer client.
    If omitted, the `region` argument of the provider is used.

* `flavor_id` - (Optional) The ID of the flavor. Exactly one of `name`, `flavor_id` is required to be set.

* `name` - (Optional) The name of the flavor. Exactly one of `name`, `flavor_id` is required to be set.

## Attributes Reference

`id` is set to the ID of the found flavor. In addition, the following attributes
are exported:

* `name` - The name of the flavor.

* `description` - The description of the flavor.

* `flavor_id` - The ID of the flavor.

* `flavor_profile_id` - The ID of the flavor profile.

* `enabled` - Is the flavor enabled.
