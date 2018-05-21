---
layout: "openstack"
page_title: "OpenStack: openstack_networking_subnet_route_v2"
sidebar_current: "docs-openstack-resource-networking-subnet-route-v2"
description: |-
  Creates a routing entry on a OpenStack V2 subnet.
---

# openstack\_networking\_subnet\_route\_v2

Creates a routing entry on a OpenStack V2 subnet.

## Example Usage

```hcl
resource "openstack_networking_router_v2" "router_1" {
  name           = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "openstack_networking_subnet_route_v2" "subnet_route_1" {
  subnet_id        = "${openstack_networking_subnet_v2.subnet_1.id}"
  destination_cidr = "10.0.1.0/24"
  next_hop         = "192.168.199.254"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to configure a routing entry on a subnet. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    routing entry.

* `subnet_id` - (Required) ID of the subnet this routing entry belongs to. Changing
    this creates a new routing entry.

* `destination_cidr` - (Required) CIDR block to match on the packet’s destination IP. Changing
    this creates a new routing entry.

* `next_hop` - (Required) IP address of the next hop gateway.  Changing
    this creates a new routing entry.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `destination_cidr` - See Argument Reference above.
* `next_hop` - See Argument Reference above.

## Notes

## Import

Routing entries can be imported using a combined ID using the following format: ``<subnet_id>-route-<destination_cidr>-<next_hop>``

```
$ terraform import openstack_networking_subnet_route_v2.subnet_route_1 686fe248-386c-4f70-9f6c-281607dad079-route-10.0.1.0/24-192.168.199.25
```
