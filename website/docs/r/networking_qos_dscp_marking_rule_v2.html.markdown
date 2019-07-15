---
layout: "openstack"
page_title: "OpenStack: openstack_networking_qos_dscp_marking_rule_v2"
sidebar_current: "docs-openstack-resource-networking-qos-dscp-marking-rule-v2"
description: |-
  Manages a V2 Neutron QoS DSCP marking rule resource within OpenStack.
---

# openstack\_networking\_qos\_dscp\_marking\_rule\_v2

Manages a V2 Neutron QoS DSCP marking rule resource within OpenStack.

## Example Usage

### Create a QoS Policy with some DSCP marking rule

```hcl
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name        = "qos_policy_1"
  description = "dscp_mark"
}

resource "openstack_networking_qos_dscp_marking_rule_v2" "dscp_marking_rule_1" {
  qos_policy_id = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  dscp_mark     = 26
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron QoS DSCP marking rule. If omitted, the
    `region` argument of the provider is used. Changing this creates a new QoS DSCP marking rule.
    
* `qos_policy_id` - (Required) The QoS policy reference. Changing this creates a new QoS DSCP marking rule.
   
* `dscp_mark` - (Required) The value of DSCP mark. Changing this updates the DSCP mark value existing
    QoS DSCP marking rule.
    
## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `qos_policy_id` - See Argument Reference above.
* `dscp_mark` - See Argument Reference above.

## Import

QoS DSCP marking rules can be imported using the `qos_policy_id/dscp_marking_rule_id` format, e.g.

```
$ terraform import openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1 d6ae28ce-fcb5-4180-aa62-d260a27e09ae/46dfb556-b92f-48ce-94c5-9a9e2140de94
```