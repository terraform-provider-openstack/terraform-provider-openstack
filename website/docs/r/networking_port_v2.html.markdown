---
layout: "openstack"
page_title: "OpenStack: openstack_networking_port_v2"
sidebar_current: "docs-openstack-resource-networking-port-v2"
description: |-
  Manages a V2 port resource within OpenStack.
---

# openstack\_networking\_port\_v2

Manages a V2 port resource within OpenStack.

## Example Usage

### Simple port

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

### Port with physical binding information

```hcl
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = "${openstack_networking_network_v2.network_1.id}"
  device_id      = "cdf70fcf-c161-4f24-9c70-96b3f5a54b71"
  device_owner   = "baremetal:none"
  admin_state_up = "true"

  binding {
    host_id   = "b080b9cf-46e0-4ce8-ad47-0fd4accc872b"
    vnic_type = "baremetal"
    profile   = <<EOF
{
  "local_link_information": [
    {
      "switch_info": "info1",
      "port_id": "Ethernet3/4",
      "switch_id": "12:34:56:78:9A:BC"
    },
    {
      "switch_info": "info2",
      "port_id": "Ethernet3/4",
      "switch_id": "12:34:56:78:9A:BD"
    }
  ],
  "vlan_type": "allowed"
}
EOF
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    port.

* `name` - (Optional) A unique name for the port. Changing this
    updates the `name` of an existing port.

* `description` - (Optional) Human-readable description of the port. Changing
    this updates the `description` of an existing port.

* `network_id` - (Required) The ID of the network to attach the port to. Changing
    this creates a new port.

* `admin_state_up` - (Optional) Administrative up/down status for the port
    (must be `true` or `false` if provided). Changing this updates the
    `admin_state_up` of an existing port.

* `mac_address` - (Optional) Specify a specific MAC address for the port. Changing
    this creates a new port.

* `tenant_id` - (Optional) The owner of the port. Required if admin wants
    to create a port for another tenant. Changing this creates a new port.

* `device_owner` - (Optional) The device owner of the port. Changing this creates
    a new port.

* `security_group_ids` - (Optional - Conflicts with `no_security_groups`) A list
    of security group IDs to apply to the port. The security groups must be
    specified by ID and not name (as opposed to how they are configured with
    the Compute Instance).

* `no_security_groups` - (Optional - Conflicts with `security_group_ids`) If set to
    `true`, then no security groups are applied to the port. If set to `false` and
    no `security_group_ids` are specified, then the port will yield to the default
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

* `extra_dhcp_option` - (Optional) An extra DHCP option that needs to be configured
    on the port. The structure is described below. Can be specified multiple
    times.

* `port_security_enabled` - (Optional) Whether to explicitly enable or disable
  port security on the port. Port Security is usually enabled by default, so
  omitting argument will usually result in a value of `true`. Setting this
  explicitly to `false` will disable port security. In order to disable port
  security, the port must not have any security groups. Valid values are `true`
  and `false`.

* `value_specs` - (Optional) Map of additional options.

* `tags` - (Optional) A set of string tags for the port.

* `binding` - (Optional) The port binding allows to specify binding information
    for the port. The structure is described below.

* `dns_name` - (Optional) The port DNS name. Available, when Neutron DNS extension
    is enabled.
    
* `qos_policy_id` - (Optional) Reference to the associated QoS policy.

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

The `extra_dhcp_option` block supports:

* `name` - (Required) Name of the DHCP option.

* `value` - (Required) Value of the DHCP option.

* `ip_version` - (Optional) IP protocol version. Defaults to 4.

The `binding` block supports:

* `host_id` - (Optional) The ID of the host to allocate port on.

* `profile` - (Optional) Custom data to be passed as `binding:profile`. Data
    must be passed as JSON.

* `vnic_type` - (Optional) VNIC type for the port. Can either be `direct`,
    `direct-physical`, `macvtap`, `normal`, `baremetal` or `virtio-forwarder`.
    Default value is `normal`.

* `vif_details` - (Computed) A map of JSON strings containing additional
    details for this specific binding.

* `vif_type` - (Computed) The VNIC type of the port binding.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `description` - See Argument Reference above.
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
* `extra_dhcp_option` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the port, which have been
  explicitly and implicitly added.
* `binding` - See Argument Reference above.
* `dns_name` - See Argument Reference above.
* `dns_assignment` - The list of maps representing port DNS assignments.
* `qos_policy_id` - See Argument Reference above.

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
