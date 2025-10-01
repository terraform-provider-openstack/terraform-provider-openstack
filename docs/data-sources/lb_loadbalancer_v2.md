---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_loadbalancer_v2"
sidebar_current: "docs-openstack-datasource-lb-loadbalancer-v2"
description: |-
  Get information on an OpenStack Load Balancer.
---

# openstack\_lb\_loadbalancer\_v2

Use this data source to get the ID of an OpenStack Load Balancer.

## Example Usage

```hcl
data "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Load Balancer client.
    If omitted, the `region` argument of the provider is used.

* `loadbalancer_id` - (Optional) The ID of the loadbalancer. Exactly one of
  `name`, `loadbalancer_id` is required to be set.

* `name` - (Optional) The name of the loadbalancer. Exactly one of `name`,
  `loadbalancer_id` is required to be set.

* `description` - (Optional) The human-readable description for the loadbalancer.

* `vip_address` - (Optional) The IP address of the loadbalancer's virtual IP (VIP).

* `tags` - (Optional) A set of tags applied to the loadbalancer. The load balancer
  will be returned if it has all of the specified tags.

* `tags_any` - (Optional) A set of tags. The load balancer will be returned if
  it has at least one of the specified tags.

* `tags_not` - (Optional) A set of tags. The load balancer will be returned if
  it does not have all of the specified tags.

* `tags_not_any` - (Optional) A set of tags. The load balancer will be returned
  if it does not have any of the specified tags.

## Attributes Reference

`id` is set to the ID of the found loadbalancer. In addition, the following
attributes are exported:

* `name` - The name of the loadbalancer.

* `description` - The description of the loadbalancer.

* `admin_state_up` - The administrative state of the loadbalancer (true/false).

* `project_id` - The owner (project/tenant) ID of the loadbalancer.

* `provisioning_status` - The provisioning status of the loadbalancer.

* `operating_status` - The operating status of the loadbalancer.

* `vip_address` - The IP address of the loadbalancer’s virtual IP (VIP).

* `vip_port_id` - The port ID associated with the VIP.

* `vip_subnet_id` - The subnet ID associated with the VIP.

* `vip_network_id` - The network ID associated with the VIP.

* `vip_qos_policy_id` - The QoS policy ID associated with the VIP, if any.

* `flavor_id` - The flavor ID used by the loadbalancer.

* `availability_zone` - The availability zone of the loadbalancer.

* `loadbalancer_provider` - The loadbalancer driver/provider used by Octavia
  (for example, `amphora`).

* `tags` - A set of tags applied to the loadbalancer.

* `listeners` - A list of listener IDs (UUIDs) associated with the loadbalancer.

* `pools` - A list of pool IDs (UUIDs) associated with the loadbalancer.

* `additional_vips` - A list of additional VIP IP addresses associated with
  the loadbalancer.
