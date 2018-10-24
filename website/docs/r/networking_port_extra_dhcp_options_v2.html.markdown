---
layout: "openstack"
page_title: "OpenStack: openstack_networking_port_extradhcp_options_v2"
sidebar_current: "docs-openstack-resource-networking-port-extradhcp-options-v2"
description: |-
  Manages a V2 port DHCP options resource within OpenStack.
---

# openstack\_networking\_port\_extradhcp\_options_v2

Manages a V2 port DHCP options resource within OpenStack.

## Example Usage

```hcl
resource "openstack_networking_network_v2" "network_1" {
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id     = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id  =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}

resource "openstack_networking_port_extradhcp_options_v2" "opts_1" {
  name    = "%s"
  port_id = "${openstack_networking_port_v2.port_1.id}"

  extra_dhcp_opts {
    opt_name  = "bootfile-name"
    opt_value = "testfile.1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to create DHCP options. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    port DHCP options resource.

* `name` - (Optional) A unique name for the port DHCP options set. Changing this
    creates a new DHCP options resource.

* `port_id` - (Required) The ID of the port to add the DHCP options to. Changing
    this creates a new port DHCP options resource.

* `extra_dhcp_opts` - (Required) An array of desired DHCP options that needs to
    be configured on the port. The structure is described below.

The `extra_dhcp_opts` block supports:

* `ip_address` - (Required) Name of the DHCP option.

* `ip_address` - (Required) Value of the DHCP option.

* `ip_version` - (Optional) IP protocol version. Defaults to 4.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `port_id` - See Argument Reference above.
* `extra_dhcp_opts` - See Argument Reference above.

## Import

DHCP options can be imported using the `port_id`, e.g.

```
$ terraform import openstack_networking_port_extradhcp_options_v2.opts_1 ed41351c-b551-4088-8a29-6a164067641c
```
