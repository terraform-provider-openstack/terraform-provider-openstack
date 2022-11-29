---
layout: "openstack"
page_title: "OpenStack: openstack_networking_secgroup_v2"
sidebar_current: "docs-openstack-resource-networking-secgroup-v2"
description: |-
  Manages a V2 Neutron security group resource within OpenStack.
---

# openstack\_networking\_secgroup\_v2

Manages a V2 neutron security group resource within OpenStack.
Unlike Nova security groups, neutron separates the group from the rules
and also allows an admin to target a specific tenant_id.

~> **NOTE on Security Groups and Security Group Rules:** We currently support the
definition of security group rules with both, openstack_networking_secgroup_v2 and
openstack_networking_secgroup_rule_v2. It is at this time not possible to mix
the usage of both resource. Doing so will lead to unpredictable behavior.

## Example Usage

```hcl
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
}
```

```hcl
resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name        = "secgroup_2"

  rule {
    port_range_min   = 22
    port_range_max   = 22
    protocol         = "tcp"
    ethertype        = "IPv4"
    direction        = "ingress"
    remote_ip_prefix = "0.0.0.0/0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    security group.

* `name` - (Required) A unique name for the security group.

* `description` - (Optional) A unique name for the security group.

* `tenant_id` - (Optional) The owner of the security group. Required if admin
    wants to create a port for another tenant. Changing this creates a new
    security group.

* `delete_default_rules` - (Optional) Whether or not to delete the default
    egress security rules. This is `false` by default. See the below note
    for more information.

* `tags` - (Optional) A set of string tags for the security group.

* `rule` - (Optional) A configuration block of security group rules. Can be specified
    multiple times for each security group rule. The structure is documented below.

The `rule` block supports:

* `description` - (Optional) A description of the rule. Changing this creates a
    new security group rule.

* `direction` - (Required) The direction of the rule, valid values are **ingress**
    or **egress**. Changing this creates a new security group rule.

* `ethertype` - (Required) The layer 3 protocol type, valid values are **IPv4**
    or **IPv6**. Changing this creates a new security group rule.

* `protocol` - (Optional) The layer 4 protocol type, valid values are following.
    This is required if you want to specify a port range. Changing this creates
    a new security group rule.
  * **tcp**
  * **udp**
  * **icmp**
  * **ah**
  * **dccp**
  * **egp**
  * **esp**
  * **gre**
  * **igmp**
  * **ipv6-encap**
  * **ipv6-frag**
  * **ipv6-icmp**
  * **ipv6-nonxt**
  * **ipv6-opts**
  * **ipv6-route**
  * **ospf**
  * **pgm**
  * **rsvp**
  * **sctp**
  * **udplite**
  * **vrrp**

* `port_range_min` - (Required) The lower part of the allowed port range, valid
    integer value needs to be between 1 and 65535. Changing this creates a new
    security group rule.

* `port_range_max` - (Required) The higher part of the allowed port range, valid
    integer value needs to be between 1 and 65535. Changing this creates a new
    security group rule.

* `remote_ip_prefix` - (Optional) The remote CIDR, the value needs to be a valid
    CIDR (i.e. 192.168.0.0/16). Changing this creates a new security group rule.

* `remote_group_id` - (Optional) The remote group id, the value needs to be an
    Openstack ID of a security group in the same tenant. Changing this creates
    a new security group rule.

* `self` - (Optional) Required if `remote_ip_prefix` and `remote_group_id` is
    empty. If true, the security group itself will be added as a source to this
    security group rule. Cannot be combined with `remote_ip_prefix` or
    `remote_group_id`.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `rule` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the security group, which have
  been explicitly and implicitly added.

## Default Security Group Rules

In most cases, OpenStack will create some egress security group rules for each
new security group. These security group rules will not be managed by
Terraform, so if you prefer to have *all* aspects of your infrastructure
managed by Terraform, set `delete_default_rules` to `true` and then create
separate security group rules such as the following:

```hcl
resource "openstack_networking_secgroup_rule_v2" "secgroup_rule_v4" {
  direction         = "egress"
  ethertype         = "IPv4"
  security_group_id = "${openstack_networking_secgroup_v2.secgroup.id}"
}

resource "openstack_networking_secgroup_rule_v2" "secgroup_rule_v6" {
  direction         = "egress"
  ethertype         = "IPv6"
  security_group_id = "${openstack_networking_secgroup_v2.secgroup.id}"
}
```

Please note that this behavior may differ depending on the configuration of
the OpenStack cloud. The above illustrates the current default Neutron
behavior. Some OpenStack clouds might provide additional rules and some might
not provide any rules at all (in which case the `delete_default_rules` setting
is moot).

## Import

Security Groups can be imported using the `id`, e.g.

```
$ terraform import openstack_networking_secgroup_v2.secgroup_1 38809219-5e8a-4852-9139-6f461c90e8bc
```
