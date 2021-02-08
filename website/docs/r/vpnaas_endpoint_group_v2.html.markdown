---
layout: "openstack"
page_title: "OpenStack: openstack_vpnaas_endpoint_group_v2"
sidebar_current: "docs-openstack-resource-vpnaas-endpoint-group-v2"
description: |-
  Manages a V2 Neutron Endpoint Group resource within OpenStack.
---

# openstack\_vpnaas\_endpoint\_group\_v2

Manages a V2 Neutron Endpoint Group resource within OpenStack.

## Example Usage

```hcl
resource "openstack_vpnaas_endpoint_group_v2" "group_1" {
  name = "Group 1"
  type = "cidr"
  endpoints = ["10.2.0.0/24",
  "10.3.0.0/24", ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create an endpoint group. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    group.

* `name` - (Optional) The name of the group. Changing this updates the name of
    the existing group.

* `tenant_id` - (Optional) The owner of the group. Required if admin wants to
    create an endpoint group for another project. Changing this creates a new group.

* `description` - (Optional) The human-readable description for the group.
    Changing this updates the description of the existing group.

* `type` -  The type of the endpoints in the group. A valid value is subnet, cidr, network, router, or vlan.
    Changing this creates a new group.
    
* `endpoints` - List of endpoints of the same type, for the endpoint group. The values will depend on the type.
    Changing this creates a new group.
    
* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `type` - See Argument Reference above.
* `endpoints` - See Argument Reference above.
* `value_specs` - See Argument Reference above.


## Import

Groups can be imported using the `id`, e.g.

```
$ terraform import openstack_vpnaas_endpoint_group_v2.group_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
