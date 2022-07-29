---
layout: "openstack"
page_title: "OpenStack: openstack_lb_member_v2"
sidebar_current: "docs-openstack-resource-lb-member-v2"
description: |-
  Manages a V2 member resource within OpenStack.
---

# openstack\_lb\_member\_v2

Manages a V2 member resource within OpenStack.

## Example Usage

```hcl
resource "openstack_lb_member_v2" "member_1" {
  pool_id       = "935685fb-a896-40f9-9ff4-ae531a3a00fe"
  address       = "192.168.199.23"
  protocol_port = 8080
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a member. If omitted, the `region`
  argument of the provider is used. Changing this creates a new member.

* `pool_id` - (Required) The id of the pool that this member will be assigned
  to. Changing this creates a new member.

* `subnet_id` - (Optional) The subnet in which to access the member. Changing
  this creates a new member.

* `name` - (Optional) Human-readable name for the member.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the member.  Only administrative users can specify a tenant UUID
  other than their own. Changing this creates a new member.

* `address` - (Required) The IP address of the member to receive traffic from
  the load balancer. Changing this creates a new member.

* `protocol_port` - (Required) The port on which to listen for client traffic.
  Changing this creates a new member.

* `weight` - (Optional)  A positive integer value that indicates the relative
  portion of traffic that this member should receive from the pool. For
  example, a member with a weight of 10 receives five times as much traffic
  as a member with a weight of 2. Defaults to 1.

* `admin_state_up` - (Optional) The administrative state of the member.
  A valid value is true (UP) or false (DOWN). Defaults to true.

* `monitor_address` - (Optional) An alternate IP address used for health monitoring a backend member.
  Available only for Octavia

* `monitor_port` - (Optional) An alternate protocol port used for health monitoring a backend member.
  Available only for Octavia

* `backup` - (Optional) Boolean that indicates whether that member works as a backup or not. Available 
  only for Octavia >= 2.1.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the member.
* `name` - See Argument Reference above.
* `weight` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `pool_id` - See Argument Reference above.
* `address` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
* `monitor_address` - See Argument reference above.
* `monitor_port` - See Argument reference above.
* `backup` - See Argument reference above.

## Import

Load Balancer Pool Member can be imported using the Pool ID and Member ID
separated by a slash, e.g.:

```
$ terraform import openstack_lb_member_v2.member_1 c22974d2-4c95-4bcb-9819-0afc5ed303d5/9563b79c-8460-47da-8a95-2711b746510f
```
