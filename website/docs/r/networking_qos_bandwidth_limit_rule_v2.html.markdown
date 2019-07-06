---
layout: "openstack"
page_title: "OpenStack: openstack_networking_qos_bandwidth_limit_rule_v2"
sidebar_current: "docs-openstack-resource-networking-qos-bandwidth-limit-rule-v2"
description: |-
  Manages a V2 Neutron QoS bandwidth limit rule resource within OpenStack.
---

# openstack\_networking\_qos\_bandwidth\_limit\_rule\_v2

Manages a V2 Neutron QoS bandwidth limit rule resource within OpenStack.

## Example Usage

### Create a QoS Policy with some bandwidth limit rule

```hcl
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name        = "qos_policy_1"
  description = "bw_limit"
}

resource "openstack_networking_qos_bandwidth_limit_rule_v2" "bw_limit_rule_1" {
  qos_policy_id  = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  max_kbps       = 3000
  max_burst_kbps = 300
  direction      = "egress"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron QoS bandwidth limit rule. If omitted, the
    `region` argument of the provider is used. Changing this creates a new QoS bandwidth limit rule.
    
* `qos_policy_id` - (Required) The QoS policy reference. Changing this creates a new QoS bandwidth limit rule.
   
* `max_kbps` - (Required) The maximum kilobits per second of a QoS bandwidth limit rule. Changing this updates the
    maximum kilobits per second of the existing QoS bandwidth limit rule.

* `max_burst_kbps` - (Optional) The maximum burst size in kilobits of a QoS bandwidth limit rule. Changing this updates the
    maximum burst size in kilobits of the existing QoS bandwidth limit rule.
   
* `direction` - (Optional) The direction of traffic. Defaults to "egress". Changing this updates the direction of the
    existing QoS bandwidth limit rule.
    
## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `qos_policy_id` - See Argument Reference above.
* `max_kbps` - See Argument Reference above.
* `max_burst_kbps` - See Argument Reference above.
* `direction` - See Argument Reference above.

## Import

QoS bandwidth limit rules can be imported using the `qos_policy_id/bandwidth_limit_rule` format, e.g.

```
$ terraform import openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1 d6ae28ce-fcb5-4180-aa62-d260a27e09ae/46dfb556-b92f-48ce-94c5-9a9e2140de94
```