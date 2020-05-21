---
layout: "openstack"
page_title: "OpenStack: openstack_networking_qos_dscp_marking_rule_v2"
sidebar_current: "docs-openstack-datasource-networking-qos-dscp-marking-rule-v2"
description: |-
  Get information on an OpenStack QoS DSCP marking rule.
---

# openstack\_networking\_qos\_dscp\_marking\_rule\_v2

Use this data source to get the ID of an available OpenStack QoS DSCP marking rule.

## Example Usage

```hcl
data "openstack_networking_qos_dscp_marking_rule_v2" "qos_dscp_marking_rule_1" {
  dscp_mark = 26
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron QoS DSCP marking rule. If omitted, the
    `region` argument of the provider is used.

* `qos_policy_id` - (Required) The QoS policy reference.

* `dscp_mark` - (Optional) The value of a DSCP mark.


## Attributes Reference

`id` is set to the `qos_policy_id/dscp_marking_rule_id` format of the found QoS DSCP marking rule.
In addition, the following attributes are exported:

* `region` - See Argument Reference above.
* `qos_policy_id` - See Argument Reference above.
* `dscp_mark` - See Argument Reference above.