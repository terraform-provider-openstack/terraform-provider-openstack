---
subcategory: "Networking / Neutron"
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
    "${openstack_fw_rule_v2.rule_1.id}",
  ]
}

resource "openstack_fw_policy_v2" "policy_2" {
  name = "firewall_egress_policy"

  rules = [
    "${openstack_fw_rule_v2.rule_2.id}",
  ]
}

resource "openstack_fw_group_v2" "group_1" {
  name      = "firewall_group"
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.policy_1.id}"
  egress_firewall_policy_id = "${openstack_fw_policy_v2.policy_2.id}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the v1 networking client.
    A networking client is needed to create a firewall. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    firewall.

* `policy_id` - (Required) The policy resource id for the firewall. Changing
    this updates the `policy_id` of an existing firewall.

* `name` - (Optional) A name for the firewall. Changing this
    updates the `name` of an existing firewall.

* `description` - (Required) A description for the firewall. Changing this
    updates the `description` of an existing firewall.

* `admin_state_up` - (Optional) Administrative up/down status for the firewall
    (must be "true" or "false" if provided - defaults to "true").
    Changing this updates the `admin_state_up` of an existing firewall.

* `tenant_id` - (Optional) The owner of the floating IP. Required if admin wants
    to create a firewall for another tenant. Changing this creates a new
    firewall.

* `associated_routers` - (Optional) Router(s) to associate this firewall instance
    with. Must be a list of strings. Changing this updates the associated routers
    of an existing firewall. Conflicts with `no_routers`.

* `no_routers` - (Optional) Should this firewall not be associated with any routers
    (must be "true" or "false" if provide - defaults to "false").
    Conflicts with `associated_routers`.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `policy_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `associated_routers` - See Argument Reference above.
* `no_routers` - See Argument Reference above.

## Import

Firewalls can be imported using the `id`, e.g.

```
$ terraform import openstack_fw_group_v2.group_1 c9e39fb2-ce20-46c8-a964-25f3898c7a97
```
