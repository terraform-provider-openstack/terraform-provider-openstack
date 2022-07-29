---
layout: "openstack"
page_title: "OpenStack: openstack_lb_members_v2"
sidebar_current: "docs-openstack-resource-lb-members-v2"
description: |-
  Manages a V2 members resource within OpenStack (batch members update).
---

# openstack\_lb\_members\_v2

Manages a V2 members resource within OpenStack (batch members update).

~> **Note:** This resource has attributes that depend on octavia minor versions.
Please ensure your Openstack cloud supports the required [minor version](../#octavia-api-versioning).

~> **Note:** This resource works only within [Octavia API](../#use_octavia). For
legacy Neutron LBaaS v2 extension please use
[openstack_lb_member_v2](lb_member_v2.html) resource.

## Example Usage

```hcl
resource "openstack_lb_members_v2" "members_1" {
  pool_id = "935685fb-a896-40f9-9ff4-ae531a3a00fe"

  member {
    address       = "192.168.199.23"
    protocol_port = 8080
  }

  member {
    address       = "192.168.199.24"
    protocol_port = 8080
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create pool members. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  members resource.

* `pool_id` - (Required) The id of the pool that members will be assigned to.
  Changing this creates a new members resource.

* `member` - (Optional) A set of dictionaries containing member parameters. The
  structure is described below.

The `member` block supports:

* `subnet_id` - (Optional) The subnet in which to access the member.

* `name` - (Optional) Human-readable name for the member.

* `address` - (Required) The IP address of the members to receive traffic from
  the load balancer.

* `protocol_port` - (Required) The port on which to listen for client traffic.

* `weight` - (Optional)  A positive integer value that indicates the relative
  portion of traffic that this members should receive from the pool. For
  example, a member with a weight of 10 receives five times as much traffic
  as a member with a weight of 2. Defaults to 1.

* `monitor_port` - (Optional) An alternate protocol port used for health 
  monitoring a backend member.

* `monitor_address` - (Optional) An alternate IP address used for health 
monitoring a backend member.

* `admin_state_up` - (Optional) The administrative state of the member.
  A valid value is true (UP) or false (DOWN). Defaults to true.

* `backup` - (Optional) A bool that indicates whether the member is
  backup. **Requires octavia minor version 2.1 or later**.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the members.
* `pool_id` - See Argument Reference above.
* `member` - See Argument Reference above.

## Import

Load Balancer Pool Members can be imported using the Pool ID, e.g.:

```
$ terraform import openstack_lb_members_v2.members_1 c22974d2-4c95-4bcb-9819-0afc5ed303d5
```
