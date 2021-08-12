---
layout: "openstack"
page_title: "OpenStack: openstack_networking_portforwarding_v2"
sidebar_current: "docs-openstack-resource-networking-portforwarding-v2"
description: |-
  Manages a V2 port forwarding resource within OpenStack.
---

# openstack\_networking\_portforwarding_v2

Manages a V2 portforwarding resource within OpenStack.

## Example Usage

### Simple portforwarding

```hcl
resource "openstack_networking_portforwarding_v2" "pf_1" {
  floatingip_id    = "7a52eb59-7d47-415d-a884-046666a6fbae"
  external_port    = 7233
  internal_port    = 25
  internal_port_id = "b930d7f6-ceb7-40a0-8b81-a425dd994ccf"
  protocol         = "tcp"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to create a port forwarding. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    port forwarding.

* `floatingip_id` - The ID of the Neutron floating IP address. Changing this creates a new port forwarding.

* `internal_port_id` - The ID of the Neutron port associated with the port forwarding. Changing
    this updates the `internal_port_id` of an existing port forwarding.

* `internal_ip_address` - The fixed IPv4 address of the Neutron port associated with the port forwarding.
    Changing this updates the `internal_ip_address` of an existing port forwarding.

* `internal_port` - The TCP/UDP/other protocol port number of the Neutron port fixed IP address associated to the
    port forwarding. Changing this updates the `internal_port` of an existing port forwarding.

* `external_port` - The TCP/UDP/other protocol port number of the port forwarding. Changing this
    updates the `external_port` of an existing port forwarding.

* `protocol` - The IP protocol used in the port forwarding. Changing this updates the `protocol`
    of an existing port forwarding.

* `tenant_id` - (Optional) The owner of the port forwarding. Required if admin wants
    to create a port forwarding for another tenant. Changing this creates a new port forwarding.

* `description` - (Optional) A text describing the port forwarding. Changing this
    updates the `description` of an existing port forwarding.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `id` - The ID of the floating IP port forwarding.
* `floatingip_id` - See Argument Reference above.
* `internal_port_id` - See Argument Reference above.
* `internal_ip_address` - See Argument Reference above.
* `internal_port` - See Argument Reference above.
* `external_port` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `description` - See Argument Reference above.
