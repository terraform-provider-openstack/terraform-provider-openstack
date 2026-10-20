---
subcategory: "Networking / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_networking_router_routes_v2"
sidebar_current: "docs-openstack-resource-networking-router-routes-v2"
description: |-
  Creates routing entries on a OpenStack V2 router.
---

# openstack\_networking\_router\_routes\_v2

Creates routing entries on a OpenStack V2 router.

~> **Note:** This resource uses the OpenStack Neutron `extraroute-atomic`
extension. If your environment does not have this extension, you should use the
`openstack_networking_router_route_v2` resource to add routes instead.

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
  network_id = openstack_networking_network_v2.network_1.id
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "openstack_networking_router_interface_v2" "int_1" {
  router_id = openstack_networking_router_v2.router_1.id
  subnet_id = openstack_networking_subnet_v2.subnet_1.id
}

resource "openstack_networking_router_routes_v2" "router_routes_1" {
  router_id = openstack_networking_router_interface_v2.int_1.router_id

  routes {
    destination_cidr = "10.0.1.0/24"
    next_hop         = "192.168.199.254"
  }
  routes {
    destination_cidr = "10.0.2.0/24"
    next_hop         = "192.168.199.254"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
  A networking client is needed to configure routing entres on a router. If
  omitted, the `region` argument of the provider is used. Changing this creates
  new routing entries.

* `router_id` - (Required) ID of the router these routing entries belong to.
  Changing this creates new routing entries.

* `routes` - (Optional) A set of routing entries to add to the router.

The `routes` block supports the following arguments:

* `destination_cidr` - (Required) CIDR block to match on the packet’s
  destination IP.

* `next_hop` - (Required) IP address of the next hop gateway.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `router_id` - See Argument Reference above.
* `routes` - See Argument Reference above.

## Notes

The `next_hop` IP address must be directly reachable from the router at the
`openstack_networking_router_routes_v2` resource creation time.  You can
ensure that by explicitly specifying a dependency on the
`openstack_networking_router_interface_v2` resource that connects the next
hop to the router, as in the example above.

## Import

Routing entries can be imported using a router `id`:

```shell
terraform import openstack_networking_router_routes_v2.router_routes_1 686fe248-386c-4f70-9f6c-281607dad079
```
