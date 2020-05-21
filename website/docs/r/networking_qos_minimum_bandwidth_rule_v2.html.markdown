---
layout: "openstack"
page_title: "OpenStack: openstack_networking_qos_minimum_bandwidth_rule_v2"
sidebar_current: "docs-openstack-resource-networking-qos-minimum-bandwidth-rule-v2"
description: |-
  Manages a V2 Neutron QoS minimum bandwidth rule resource within OpenStack.
---

# openstack\_networking\_qos\_minimum\_bandwidth\_rule\_v2

Manages a V2 Neutron QoS minimum bandwidth rule resource within OpenStack.

## Example Usage

### Create a QoS Policy with some minimum bandwidth rule

```hcl
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name        = "qos_policy_1"
  description = "min_kbps"
}

resource "openstack_networking_qos_minimum_bandwidth_rule_v2" "minimum_bandwidth_rule_1" {
  qos_policy_id = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  min_kbps      = 200
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron QoS minimum bandwidth rule. If omitted, the
    `region` argument of the provider is used. Changing this creates a new QoS minimum bandwidth rule.
    
* `qos_policy_id` - (Required) The QoS policy reference. Changing this creates a new QoS minimum bandwidth rule.
   
* `min_kbps` - (Required) The minimum kilobits per second. Changing this updates the min kbps value of the existing
    QoS minimum bandwidth rule.

* `direction` - (Optional) The direction of traffic. Defaults to "egress". Changing this updates the direction of the
    existing QoS minimum bandwidth rule.
    
## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `qos_policy_id` - See Argument Reference above.
* `min_kbps` - See Argument Reference above.
* `direction` - See Argument Reference above.

## Import

QoS minimum bandwidth rules can be imported using the `qos_policy_id/minimum_bandwidth_rule_id` format, e.g.

```
$ terraform import openstack_networking_qos_minimum_bandwidth_rule_v2.minimum_bandwidth_rule_1 d6ae28ce-fcb5-4180-aa62-d260a27e09ae/46dfb556-b92f-48ce-94c5-9a9e2140de94
```