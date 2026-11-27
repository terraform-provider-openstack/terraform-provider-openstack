---
subcategory: "Networking / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_networking_address_group_v2"
sidebar_current: "docs-openstack-resource-networking-address-group-v2"
description: |-
  Manages a V2 Neutron address group resource within OpenStack.
---

# openstack\_networking\_address\_group\_v2

Manages a V2 neutron address group resource within OpenStack.

## Example Usage

```hcl
resource "openstack_networking_address_group_v2" "group_1" {
  name        = "group_1"
  description = "My neutron address group"
  addresses = [
    "192.168.0.1/32",
    "2001:db8::1/128",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
  If omitted, the `region` argument of the provider is used. Changing this
  creates a new address group.

* `name` - (Optional) A name of the address group.

* `description` - (Optional) A description of the address group.

* `project_id` - (Optional) The owner of the address group. Required if admin
  wants to create a group for a specific project. Changing this creates a new
  address group.

* `addresses` - (Required) A list of CIDR blocks that define the addresses in
  the address group. Each address must be a valid IPv4 or IPv6 CIDR block.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `addresses` - See Argument Reference above.

## Import

Address Groups can be imported using the `id`, e.g.

```shell
terraform import openstack_networking_address_group_v2.group_1 782fef9c-d03c-400a-9735-2f9af5681cb3
```
