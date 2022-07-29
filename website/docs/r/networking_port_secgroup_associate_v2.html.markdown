---
layout: "openstack"
page_title: "OpenStack: openstack_networking_port_secgroup_associate_v2"
sidebar_current: "docs-openstack-resource-networking-port-secgroup-associate-v2"
description: |-
  Manages a V2 port's security groups within OpenStack.
---

# openstack\_networking\_port\_secgroup\_associate\_v2

Manages a V2 port's security groups within OpenStack. Useful, when the port was
created not by Terraform (e.g. Manila or LBaaS). It should not be used, when the
port was created directly within Terraform.

When the resource is deleted, Terraform doesn't delete the port, but unsets the
list of user defined security group IDs.  However, if `enforce` is set to `true`
and the resource is deleted, Terraform will remove all assigned security group
IDs.

## Example Usage

### Append a security group to an existing port

```hcl
data "openstack_networking_port_v2" "system_port" {
  fixed_ip = "10.0.0.10"
}

data "openstack_networking_secgroup_v2" "secgroup" {
  name = "secgroup"
}

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.system_port.id}"
  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.secgroup.id}",
  ]
}
```

### Enforce a security group to an existing port

```hcl
data "openstack_networking_port_v2" "system_port" {
  fixed_ip = "10.0.0.10"
}

data "openstack_networking_secgroup_v2" "secgroup" {
  name = "secgroup"
}

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.system_port.id}"
  enforce = "true"
  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.secgroup.id}",
  ]
}
```

### Remove all security groups from an existing port

```hcl
data "openstack_networking_port_v2" "system_port" {
  fixed_ip = "10.0.0.10"
}

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id            = "${data.openstack_networking_port_v2.system_port.id}"
  enforce            = "true"
  security_group_ids = []
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to manage a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    resource.

* `port_id` - (Required) An UUID of the port to apply security groups to.

* `security_group_ids` - (Required) A list of security group IDs to apply to
    the port. The security groups must be specified by ID and not name (as
    opposed to how they are configured with the Compute Instance).

* `enforce` - (Optional) Whether to replace or append the list of security
    groups, specified in the `security_group_ids`. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `port_id` - See Argument Reference above.
* `security_group_ids` - See Argument Reference above.
* `all_security_group_ids` - The collection of Security Group IDs on the port
  which have been explicitly and implicitly added.

## Import

Port security group association can be imported using the `id` of the port, e.g.

```
$ terraform import openstack_networking_port_secgroup_associate_v2.port_1 eae26a3e-1c33-4cc1-9c31-0cd729c438a1
```
