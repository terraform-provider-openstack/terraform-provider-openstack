---
layout: "openstack"
page_title: "OpenStack: openstack_networking_port_v2"
sidebar_current: "docs-openstack-resource-networking-port-v2"
description: |-
  Manages a V2 port resource within OpenStack.
---

# openstack\_networking\_port_v2

Manages a V2 port resource within OpenStack.

## Example Usage

```hcl
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    port.

* `name` - (Optional) A unique name for the port. Changing this
    updates the `name` of an existing port.

* `network_id` - (Required) The ID of the network to attach the port to. Changing
    this creates a new port.

* `admin_state_up` - (Optional) Administrative up/down status for the port
    (must be "true" or "false" if provided). Changing this updates the
    `admin_state_up` of an existing port.

* `mac_address` - (Optional) Specify a specific MAC address for the port. Changing
    this creates a new port.

* `tenant_id` - (Optional) The owner of the Port. Required if admin wants
    to create a port for another tenant. Changing this creates a new port.

* `device_owner` - (Optional) The device owner of the Port. Changing this creates
    a new port.

* `security_group_ids` - (Optional - Conflicts with `no_security_groups`) A list
    of security group IDs to apply to the port. The security groups must be
    specified by ID and not name (as opposed to how they are configured with
    the Compute Instance).

* `no_security_groups` - (Optional - Conflicts with `security_group_ids`) If set to
    `true`, then no security groups are applied to the port. If set to `false` and
    no `security_group_ids` are specified, then the Port will yield to the default
    behavior of the Networking service, which is to usually apply the "default"
    security group.

* `device_id` - (Optional) The ID of the device attached to the port. Changing this
    creates a new port.

* `fixed_ip` - (Optional - Conflicts with `no_fixed_ip`) An array of desired IPs for
    this port. The structure is described below.

* `no_fixed_ip` - (Optional - Conflicts with `fixed_ip`) Create a port with no fixed
    IP address. This will also remove any fixed IPs previously set on a port. `true`
    is the only valid value for this argument.

* `allowed_address_pairs` - (Optional) An IP/MAC Address pair of additional IP
    addresses that can be active on this port. The structure is described
    below.

* `extra_dhcp_opts` - (Required) An array of desired DHCP options that needs to
    be configured on the port. The structure is described below.

* `value_specs` - (Optional) Map of additional options.

The `fixed_ip` block supports:

* `subnet_id` - (Required) Subnet in which to allocate IP address for
this port.

* `ip_address` - (Optional) IP address desired in the subnet for this port. If
you don't specify `ip_address`, an available IP address from the specified
subnet will be allocated to this port. This field will not be populated if it
is left blank or omitted. To retrieve the assigned IP address, use the
`all_fixed_ips` attribute.

The `allowed_address_pairs` block supports:

* `ip_address` - (Required) The additional IP address.

* `mac_address` - (Optional) The additional MAC address.

* `tags` - (Optional) A set of string tags for the port.

The `extra_dhcp_opts` block supports:

* `ip_address` - (Required) Name of the DHCP option.

* `ip_address` - (Required) Value of the DHCP option.

* `ip_version` - (Optional) IP protocol version. Defaults to 4.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `mac_address` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `device_owner` - See Argument Reference above.
* `security_group_ids` - See Argument Reference above.
* `device_id` - See Argument Reference above.
* `fixed_ip` - See Argument Reference above.
* `all_fixed_ips` - The collection of Fixed IP addresses on the port in the
  order returned by the Network v2 API.
* `all_security_group_ids` - The collection of Security Group IDs on the port
  which have been explicitly and implicitly added.
* `extra_dhcp_opts` - See Argument Reference above.
* `tags` - See Argument Reference above.

## Import

Ports can be imported using the `id`, e.g.

```
$ terraform import openstack_networking_port_v2.port_1 eae26a3e-1c33-4cc1-9c31-0cd729c438a1
```

## Notes

### Ports and Instances

There are some notes to consider when connecting Instances to networks using
Ports. Please see the `openstack_compute_instance_v2` documentation for further
documentation.
