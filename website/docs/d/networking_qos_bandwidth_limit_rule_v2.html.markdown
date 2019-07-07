---
layout: "openstack"
page_title: "OpenStack: openstack_networking_qos_bandwidth_limit_rule_v2"
sidebar_current: "docs-openstack-datasource-networking-qos-bandwidth-limit-rule-v2"
description: |-
  Get information on an OpenStack QoS Bandwidth limit rule.
---

# openstack\_networking\_qos\_bandwidth\_limit\_rule\_v2

Use this data source to get the ID of an available OpenStack QoS bandwidth limit rule.

## Example Usage

```hcl
data "openstack_networking_qos_bandwidth_limit_rule_v2" "qos_bandwidth_limit_rule_1" {
  max_kbps = 300
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron QoS bandwidth limit rule. If omitted, the
    `region` argument of the provider is used.
    
* `qos_policy_id` - (Required) The QoS policy reference.
   
* `max_kbps` - (Optional) The maximum kilobits per second of a QoS bandwidth limit rule.

* `max_burst_kbps` - (Optional) The maximum burst size in kilobits of a QoS bandwidth limit rule.
   
* `direction` - (Optional) The direction of traffic.


## Attributes Reference

`id` is set to the `qos_policy_id/bandwidth_limit_rule_id` format of the found QoS bandwidth limit rule.
In addition, the following attributes are exported:

* `region` - See Argument Reference above.
* `qos_policy_id` - See Argument Reference above.
* `max_kbps` - See Argument Reference above.
* `max_burst_kbps` - See Argument Reference above.
* `direction` - See Argument Reference above.
