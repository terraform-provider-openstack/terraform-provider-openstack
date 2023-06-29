---
subcategory: "FWaaS / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_fw_group_v2"
sidebar_current: "docs-openstack-resource-fw-group-v2"
description: |-
  Manages a v2 firewall group resource within OpenStack.
---

# openstack\_fw\_group\_v2

Manages a v2 firewall group resource within OpenStack.

~> **Note:** Firewall v2 has no support for OVN currently.

## Example Usage

```hcl
resource "openstack_fw_rule_v2" "rule_1" {
  name             = "firewall_rule_2"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_port = "23"
  enabled          = "true"
}

resource "openstack_fw_rule_v2" "rule_2" {
  name             = "firewall_rule_1"
  description      = "drop NTP traffic"
  action           = "deny"
  protocol         = "udp"
  destination_port = "123"
  enabled          = "false"
}

resource "openstack_fw_policy_v2" "policy_1" {
  name = "firewall_ingress_policy"

  rules = [
    openstack_fw_rule_v2.rule_1.id,
  ]
}

resource "openstack_fw_policy_v2" "policy_2" {
  name = "firewall_egress_policy"

  rules = [
    openstack_fw_rule_v2.rule_2.id,
  ]
}

resource "openstack_fw_group_v2" "group_1" {
  name      = "firewall_group"
  ingress_firewall_policy_id = openstack_fw_policy_v2.policy_1.id
  egress_firewall_policy_id = openstack_fw_policy_v2.policy_2.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the v2 networking client.
    A networking client is needed to create a firewall group. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    firewall group.

* `name` - (Optional) A name for the firewall group. Changing this
    updates the `name` of an existing firewall.

* `description` - (Optional) A description for the firewall group. Changing this
    updates the `description` of an existing firewall group.

* `tenant_id` - (Optional) - This argument conflicts and is interchangeable with
    `project_id`. The owner of the firewall group. Required if admin wants to
    create a firewall group for another tenant. Changing this creates a new
    firewall group.

* `project_id` - (Optional) - This argument conflicts and  is interchangeable
    with `tenant_id`. The owner of the firewall group. Required if admin wants
    to create a firewall group for another project. Changing this creates a new
    firewall group.

* `ingress_firewall_policy_id` - (Optional) The ingress firewall policy resource
    id for the firewall group. Changing this updates the
    `ingress_firewall_policy_id` of an existing firewall group.

* `egress_firewall_policy_id` - (Optional) The egress firewall policy resource
    id for the firewall group. Changing this updates the
    `egress_firewall_policy_id` of an existing firewall group.

* `admin_state_up` - (Optional) Administrative up/down status for the firewall
    group (must be "true" or "false" if provided - defaults to "true").
    Changing this updates the `admin_state_up` of an existing firewall group.

* `ports` - (Optional) Port(s) to associate this firewall group
    with. Must be a list of strings. Changing this updates the associated ports
    of an existing firewall group.

* `shared` - (Optional) Sharing status of the firewall group (must be "true"
    or "false" if provided). If this is "true" the firewall group is visible to,
    and can be used in, firewalls in other tenants. Changing this updates the
    `shared` status of an existing firewall group. Only administrative users
    can specify if the firewall group should be shared.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `ingress_firewall_policy_id` - See Argument Reference above.
* `egress_firewall_policy_id` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `ports` - See Argument Reference above.
* `shared` - See Argument Reference above.
* `status` - The status of the firewall group.

## Import

Firewall groups can be imported using the `id`, e.g.

```
$ terraform import openstack_fw_group_v2.group_1 c9e39fb2-ce20-46c8-a964-25f3898c7a97
```
