---
layout: "openstack"
page_title: "OpenStack: openstack_networking_trunk_v2"
sidebar_current: "docs-openstack-resource-networking-trunk-v2"
description: |-
  Manages a networking V2 trunk resource within OpenStack.
---

# openstack\_networking\_trunk\_v2

Manages a networking V2 trunk resource within OpenStack.

## Example Usage

```hcl
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name        = "subnet_1"
  network_id  = "${openstack_networking_network_v2.network_1.id}"
  cidr        = "192.168.1.0/24"
  ip_version  = 4
  enable_dhcp = true
  no_gateway  = true
}

resource "openstack_networking_port_v2" "parent_port_1" {
  depends_on = [
    "openstack_networking_subnet_v2.subnet_1",
  ]

  name           = "parent_port_1"
  network_id     = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}

resource "openstack_networking_port_v2" "subport_1" {
  depends_on = [
    "openstack_networking_subnet_v2.subnet_1",
  ]

  name           = "subport_1"
  network_id     = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name           = "trunk_1"
  admin_state_up = "true"
  port_id        = "${openstack_networking_port_v2.parent_port_1.id}"

  sub_port {
    port_id           = "${openstack_networking_port_v2.subport_1.id}"
    segmentation_id   = 1
    segmentation_type = "vlan"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]

  network {
    port = "${openstack_networking_trunk_v2.trunk_1.port_id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to create a trunk. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    trunk.

* `name` - (Optional) A unique name for the trunk. Changing this
    updates the `name` of an existing trunk.

* `description` - (Optional) Human-readable description of the trunk. Changing this
    updates the name of the existing trunk.

* `port_id` - (Required) The ID of the port to be used as the parent port of the
    trunk. This is the port that should be used as the compute instance network
    port. Changing this creates a new trunk.

* `admin_state_up` - (Optional) Administrative up/down status for the trunk
    (must be "true" or "false" if provided). Changing this updates the
    `admin_state_up` of an existing trunk.

* `tenant_id` - (Optional) The owner of the Trunk. Required if admin wants
    to create a trunk on behalf of another tenant. Changing this creates a new trunk.

* `sub_port` - (Optional) The set of ports that will be made subports of the trunk.
    The structure of each subport is described below.

* `tags` - (Optional) A set of string tags for the port.

The `sub_port` block supports:

* `port_id` - (Required) The ID of the port to be made a subport of the trunk.

* `segmentation_type` - (Required) The segmentation technology to use, e.g., "vlan".

* `segmentation_id` - (Required) The numeric id of the subport segment.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `port_id` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `sub_port` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the trunk, which have been
  explicitly and implicitly added.
