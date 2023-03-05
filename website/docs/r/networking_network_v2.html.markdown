---
layout: "openstack"
page_title: "OpenStack: openstack_networking_network_v2"
sidebar_current: "docs-openstack-resource-networking-network-v2"
description: |-
  Manages a V2 Neutron network resource within OpenStack.
---

# openstack\_networking\_network\_v2

Manages a V2 Neutron network resource within OpenStack.

## Example Usage

```hcl
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "openstack_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "openstack_networking_port_v2" "port_1" {
  name               = "port_1"
  network_id         = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up     = "true"
  security_group_ids = ["${openstack_compute_secgroup_v2.secgroup_1.id}"]

  fixed_ip {
    subnet_id  = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.10"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["${openstack_compute_secgroup_v2.secgroup_1.name}"]

  network {
    port = "${openstack_networking_port_v2.port_1.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron network. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    network.

* `name` - (Optional) The name of the network. Changing this updates the name of
    the existing network.

* `description` - (Optional) Human-readable description of the network. Changing this
    updates the name of the existing network.

* `shared` - (Optional) Specifies whether the network resource can be accessed
    by any tenant or not. Changing this updates the sharing capabilities of the
    existing network.

* `external` - (Optional) Specifies whether the network resource has the
    external routing facility. Valid values are true and false. Defaults to
    false. Changing this updates the external attribute of the existing network.

* `tenant_id` - (Optional) The owner of the network. Required if admin wants to
    create a network for another tenant. Changing this creates a new network.

* `admin_state_up` - (Optional) The administrative state of the network.
    Acceptable values are "true" and "false". Changing this value updates the
    state of the existing network.

* `segments` - (Optional) An array of one or more provider segment objects.
  Note: most Networking plug-ins (e.g. ML2 Plugin) and drivers do not support
  updating any provider related segments attributes. Check your plug-in whether
  it supports updating.

* `value_specs` - (Optional) Map of additional options.

* `availability_zone_hints` -  (Optional) An availability zone is used to make
    network resources highly available. Used for resources with high availability
    so that they are scheduled on different availability zones. Changing this
    creates a new network.

* `tags` - (Optional) A set of string tags for the network.

* `transparent_vlan` - (Optional) Specifies whether the network resource has the
  VLAN transparent attribute set. Valid values are true and false. Defaults to
  false. Changing this updates the `transparent_vlan` attribute of the existing
  network.

* `port_security_enabled` - (Optional) Whether to explicitly enable or disable
  port security on the network. Port Security is usually enabled by default, so
  omitting this argument will usually result in a value of "true". Setting this
  explicitly to `false` will disable port security. Valid values are `true` and
  `false`.

* `mtu` - (Optional) The network MTU. Available for read-only, when Neutron
   `net-mtu` extension is enabled. Available for the modification, when
   Neutron `net-mtu-writable` extension is enabled.

* `dns_domain` - (Optional) The network DNS domain. Available, when Neutron DNS
    extension is enabled. The `dns_domain` of a network in conjunction with the
    `dns_name` attribute of its ports will be published in an external DNS
    service when Neutron is configured to integrate with such a service.
    
* `qos_policy_id` - (Optional) Reference to the associated QoS policy.

The `segments` block supports:

* `physical_network` - The physical network where this network is implemented.
* `segmentation_id` - An isolated segment on the physical network.
* `network_type` - The type of physical network.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `shared` - See Argument Reference above.
* `external` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `availability_zone_hints` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the network, which have been
  explicitly and implicitly added.
* `transparent_vlan` - See Argument Reference above.
* `segments` - An array of one or more provider segment objects.
* `port_security_enabled` - See Argument Reference above.
* `mtu` - See Argument Reference above.
* `dns_domain` - See Argument Reference above.
* `qos_policy_id` - See Argument Reference above.

## Import

Networks can be imported using the `id`, e.g.

```
$ terraform import openstack_networking_network_v2.network_1 d90ce693-5ccf-4136-a0ed-152ce412b6b9
```
